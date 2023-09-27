package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	tg "github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/url"
	"unicode"
)

type WafflerServicer interface {
	Search(search *message.Search) (*message.SearchResponse, error)
	InfoSource(urlSearch string) (*message.InfoRequest, error)
	Score(request *message.ScoreRequest) (*message.ScoreResponse, error)
	ParseSource(search *message.ParserRequest) error
}

type WafflerService struct {
	storage storage.WafflerStorager
	log     *zap.Logger
	tg      *tg.Telegram
}

func NewWafflerService(storage storage.WafflerStorager, components *component.Components) *WafflerService {
	return &WafflerService{storage: storage, log: components.Logger, tg: components.Tg}
}

var orderRecords = map[string]string{"record_text": "record_text ASC", "record_text_desc": "record_text DESC", "score": "score ASC", "score_desc": "score DESC",
	"time": "created_at ASC", "time_desc": "created_at DESC"}

func (u *WafflerService) Score(request *message.ScoreRequest) (*message.ScoreResponse, error) {
	scoreResponse := &message.ScoreResponse{
		Cursor: &message.Cursor{
			Offset: request.Cursor.Offset,
		},
		Records: []message.Record{},
	}
	records, err := u.storage.
		SelectRecordsSourceIDOffsetLimit(request.SourceId, request.Type, orderRecords[request.Order], request.Cursor.Offset, request.Limit)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		scoreResponse.Cursor = nil

	} else {
		scoreResponse.Cursor.Offset += len(records)

		scoreResponse.Records = make([]message.Record, 0, len(records))
		for i := range records {
			scoreResponse.Records = append(scoreResponse.Records, message.Record{
				RecordText: records[i].RecordText,
				Score:      records[i].Score,
				Timestamp:  records[i].CreatedAt,
			})
		}
	}

	return scoreResponse, nil
}

func (u *WafflerService) InfoSource(urlSearch string) (*message.InfoRequest, error) {
	res := &message.InfoRequest{}

	urlParse, err := url.Parse(urlSearch)
	if err != nil {
		return nil, err
	}

	switch urlParse.Host {
	case "t.me":
		channel, err := u.tg.ContactSearch(urlSearch)
		if err != nil {
			u.log.Error("error: search", zap.Error(err))
			return nil, err
		}

		res.Name = channel.Title
		res.Type = models.Telegram
	}

	return res, nil
}

func (s *WafflerService) ParseSource(search *message.ParserRequest) error {
	dataTelegram, err := s.tg.ParseChatTelegram(search.SourceURL, 20) // TODO: search.Limit
	if err != nil {
		s.log.Error("search", zap.Error(err))
	}
	chatGPT := gpt.NewChatGPT(search.Parser.Token, s.log)

	records := make([]*models.RecordDTO, 0, 1)
	source, err := s.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return err
	}
	if source != nil {
		records, err = s.storage.SelectRecordsSourceID(source.ID)

		if err != nil {
			s.log.Error("error: search", zap.Error(err))
			return err
		}
	}
	g := errgroup.Group{}
	g.SetLimit(20)
	var indexNewRecords = -1
	newRecords := make([]*models.RecordDTO, 0, len(dataTelegram.Records))
	for _, r := range dataTelegram.Records {
		indexNewRecords++
		if !containsAlphabet(r.RecordText) {
			continue
		}
		tempIndexRecords := indexNewRecords

		existText := false
		if len(records) != 0 {
			for _, record := range records { //  ToDo: optimize
				if record.RecordText == r.RecordText {
					r.RecordText = record.RecordText
					existText = true
					break
				}

			}
		}
		if existText {
			continue
		}

		g.Go(func() error {
			var err error
			res, err := chatGPT.ConstructQuestionGPT(r.RecordText, search.ScoreType)
			if res != nil {
				err = nil
				dataTelegram.Records[tempIndexRecords].Score = *res
				newRecords = append(newRecords, dataTelegram.Records[tempIndexRecords])
				return nil
			}
			if err != nil {
				s.log.Error("error: search", zap.Error(err))
				return err
			}

			return nil
		})

	}
	err = g.Wait()
	if err != nil {
		s.log.Warn("error: search", zap.Error(err))
		return err
	}

	dataTelegram.Records = newRecords

	if source == nil {
		err := s.storage.CreateSource(dataTelegram.Source)
		if err != nil {
			s.log.Error("error: create", zap.Error(err))
			return err
		}

		source, err = s.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
		if err != nil {
			s.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	if len(dataTelegram.Records) != 0 {
		for i := range dataTelegram.Records {
			dataTelegram.Records[i].SourceID = source.ID

		}

		err := s.storage.CreateRecords(dataTelegram.Records)
		if err != nil {
			s.log.Error("error: create", zap.Error(err))
			return err
		}
		records, err = s.storage.SelectRecordsSourceID(source.ID)
		if err != nil {
			s.log.Error("error: search", zap.Error(err))
			return err
		}

		records, err = s.storage.SelectRecordsSourceID(source.ID)
		if err != nil {
			s.log.Error("error: search", zap.Error(err))
			return err
		}

		updateScoreRecods(records, source, dataTelegram)

		err = s.storage.UpdateSource(source)
		if err != nil {
			s.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	return err
}

var orderSources = map[string]string{"name": "name ASC", "name_desc": "name DESC", "source": "source_type ASC", "source_desc": "source_type DESC",
	"waffler": "waffler_score ASC", "waffler_desc": "waffler_score DESC", "racizm": "racism_score ASC", "racizm_desc": "racism_score DESC"}

//  ToDo: racizm rename to racism

func (s *WafflerService) Search(search *message.Search) (*message.SearchResponse, error) {
	res := &message.SearchResponse{
		Sources: make([]models.SourceDTO, 0, search.Limit),
	}
	var err error

	if search.Cursor.Partition == 0 {
		res.Sources, err = s.storage.
			SearchLikeBySourceName(search.QueryForName, search.SourceType, search.Cursor.Offset, orderSources[search.Order], search.Limit)

		if err != nil {
			return nil, err
		}

		if len(res.Sources) == search.Limit {
			search.Cursor.Offset += len(res.Sources)
		} else {
			search.Cursor.Partition = 1
			search.Limit -= len(res.Sources)
			search.Cursor.Offset = 0
		}

	}

	if search.Cursor.Partition == 1 {
		resURL, err := s.storage.
			SearchLikeBySourceURLNotName(search.QueryForName, search.SourceType, search.Cursor.Offset, orderSources[search.Order], search.Limit)
		if err != nil {
			return nil, err
		}

		res.Sources = append(res.Sources, resURL...)

		if len(resURL) == search.Limit {
			search.Cursor.Offset += len(resURL)
		} else {
			search.Cursor = nil
		}
	}

	res.Cursor = search.Cursor

	return res, nil
}

func containsAlphabet(text string) bool {
	for _, r := range []rune(text) {
		isValid := unicode.IsLetter(r)
		if isValid {
			return true
		}
	}

	return false

}
