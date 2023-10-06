package service

import (
	"net/url"
	"unicode"

	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/language_model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	tg      data_source.DataSourcer
}

func NewWafflerService(storage storage.WafflerStorager, components *component.Components) *WafflerService {
	return &WafflerService{storage: storage, log: components.Logger, tg: components.Tg}
}

var orderRecords = map[string]string{"score": "score ASC", "score_desc": "score DESC",
	"time": "created_at ASC", "time_desc": "created_at DESC"}

func (u *WafflerService) Score(request *message.ScoreRequest) (*message.ScoreResponse, error) {
	scoreResponse := &message.ScoreResponse{
		Cursor: &message.Cursor{
			Offset: request.Cursor.Offset,
		},
		Records: []message.Record{},
	}
	records, err := u.storage.ListRecordsSourceID(request.SourceId)
	if err != nil {
		return nil, err
	}
	var racismRecords []models.RacismDTO
	switch request.Type {
	case models.Racism:
		racismRecords, err = u.storage.ListRacismRecordsSourceIDCursor(request.SourceId, orderRecords[request.Order], request.Cursor.Offset, request.Limit)
		if err != nil {
			return nil, err
		}

	}

	if len(racismRecords) == 0 {
		scoreResponse.Cursor = nil

	} else {
		scoreResponse.Cursor.Offset += len(records)

		scoreResponse.Records = make([]message.Record, 0, len(records))
		for i := range racismRecords {
			scoreResponse.Records = append(scoreResponse.Records, message.Record{
				RecordText: records[i].RecordText,
				Score:      racismRecords[i].Score,
				Timestamp:  racismRecords[i].CreatedTs,
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

func (w *WafflerService) parseSourceTypeRacism(search *message.ParserRequest, dataTelegram *data_source.DataTelegram) error {
	g := errgroup.Group{}
	g.SetLimit(20)

	newRacismRecords := make([]models.RacismDTO, 0, len(dataTelegram.Records))
	lanModel := language_model.NewChatGPTWrapper(search.Parser.Token, w.log)
	for _, r := range dataTelegram.Records {
		tempRecord := r

		text := r.RecordText
		g.Go(func() error {
			var err error
			res, err := lanModel.ConstructQuestionGPT(tempRecord.RecordText, search)
			if res != nil {
				newRacismRecords = append(newRacismRecords, models.RacismDTO{
					Score:      *res,
					ParserType: models.GPT3_5TURBO,
					CreatedTs:  tempRecord.CreatedTs,
					RecordID:   tempRecord.ID,
					SourceID:   dataTelegram.Source.ID,
				})
				return nil
			} else {
				if err != nil {
					w.log.Error("error: search", zap.Error(err))
					return err
				}
			}
			if res != nil {
				dataTelegram.Records[tempIndexRecords].Score = *res
				newRecords = append(newRecords, dataTelegram.Records[tempIndexRecords])
			}
			return nil
		})

	}
	err := g.Wait()
	if err != nil {
		w.log.Warn("error: search", zap.Error(err))
		return err
	}

	source, err := w.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	err = w.storage.CreateRacismRecords(newRacismRecords)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	if len(dataTelegram.Records) != 0 {

		racismRecords, err := w.storage.ListRacismRecords(&models.RacismDTO{SourceID: dataTelegram.Source.ID})
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}

		source.RacismScore = func() int {
			score := 0
			for _, r := range racismRecords {
				score += r.Score
			}
			return score / len(racismRecords)
		}()

		err = w.storage.UpdateSource(source)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	return err
}

func (w *WafflerService) ParseSource(search *message.ParserRequest) error {
	dataTelegram, err := w.tg.ParseChatTelegram(search.SourceURL, 20) // TODO: search.Limit
	if err != nil {
		w.log.Error("search", zap.Error(err))
	}

	var source *models.SourceDTO
	var records []models.RecordDTO

	source, err = w.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}
	if source.Name != "" {
		records, err = w.storage.ListRecordsSourceID(source.ID)

		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}
	} else {
		err = w.storage.CreateSource(dataTelegram.Source)
		if err != nil {
			w.log.Error("error: create", zap.Error(err))
			return err
		}
    
		source, err = w.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}

	}
	for i := range dataTelegram.Records {
		dataTelegram.Records[i].SourceID = source.ID
	}

	dataTelegram.Source = source

	err = w.storage.CreateRecords(dataTelegram.Records)
	if err != nil {
		w.log.Error("error: create", zap.Error(err))
		return err
	}

	var recordsFilter = make([]models.RecordDTO, 0, len(dataTelegram.Records))

	for _, r := range dataTelegram.Records {
		recordFilter := r
		if !containsAlphabet(r.RecordText) {
			continue
		}

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
		recordsFilter = append(recordsFilter, recordFilter)
	}
	dataTelegram.Records = recordsFilter

	switch search.ScoreType {
	case models.Racism:
		err = w.parseSourceTypeRacism(search, dataTelegram)

	case models.Waffler:

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

// nolint
func containsAlphabet(text string) bool {
	for _, r := range text {
		isValid := unicode.IsLetter(r)
		if isValid {
			return true
		}
	}

	return false

}
