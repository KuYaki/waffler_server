package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	tg "github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"go.uber.org/zap"
	"sort"
)

type Waffler interface {
	Search(search *message.Search) (*message.SearchResponse, error)
	InfoSource(domain string) *message.InfoRequest
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

func (u *WafflerService) Score(request *message.ScoreRequest) (*message.ScoreResponse, error) {
	records, err := u.storage.SelectRecords(request.SourceId)
	if err != nil {
		return nil, err
	}

	recordsNew := make([]message.Record, 0, request.Limit)
	for i := request.Cursor; i < len(records) && i < request.Limit+request.Cursor; i++ {
		recordsNew = append(recordsNew, message.Record{
			RecordText: records[i].RecordText,
			Score:      records[i].Score,
			Timestamp:  records[i].CreatedAt,
		})
	}

	scoreResponse := &message.ScoreResponse{}

	scoreResponse.Records = sortRecords(recordsNew, request.Order)

	scoreResponse.Cursor += len(scoreResponse.Records)

	return scoreResponse, nil
}

func sortRecords(records []message.Record, order string) []message.Record {
	var orderRecords = []string{"record_text_up", "record_text_down", "score_up", "score_down",
		"time_up", "time_down"}
	switch order {
	case orderRecords[0]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].RecordText < records[j].RecordText
		})
	case orderRecords[1]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].RecordText > records[j].RecordText
		})
	case orderRecords[2]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Score < records[j].Score
		})
	case orderRecords[3]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Score > records[j].Score
		})
	case orderRecords[4]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Timestamp.Before(records[j].Timestamp)
		})
	case orderRecords[5]:
		sort.Slice(records, func(i, j int) bool {
			return records[i].Timestamp.After(records[j].Timestamp)
		})

	}

	return records
}

func (u *WafflerService) InfoSource(domain string) *message.InfoRequest {
	res := &message.InfoRequest{}
	switch domain {

	case "t.me":
		res.Name = "Telegram"
		res.Type = models.Telegram
		return res
	}
	return nil
}

func (s *WafflerService) ParseSource(search *message.ParserRequest) error {
	dataTelegram, err := s.tg.ParseChat(search.SourceURL, 10) // TODO: search.Limit
	if err != nil {
		s.log.Error("search", zap.Error(err))
	}
	newRecords := []models.RecordDTO{}
	chatGPT := gpt.NewChatGPT(search.Parser.Token)

	for i, r := range dataTelegram.Records {
		if r.RecordText == "" {
			continue
		}

		dataTelegram.Records[i].Score, err = chatGPT.ConstructQuestionGPT(r.RecordText, search.ScoreType)
		if err != nil {
			s.log.Error("error: question gpt", zap.Error(err))
		}
		newRecords = append(newRecords, dataTelegram.Records[i])
	}

	dataTelegram.Records = newRecords

	err = s.storage.CreateSourceAndRecords(dataTelegram)
	if err != nil {
		s.log.Error("error: create", zap.Error(err))
		return err
	}

	source, err := s.storage.SearchBySourceName(dataTelegram.Source.Name)
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return err
	}
	records, err := s.storage.SelectRecords(source[0].ID)
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return err
	}

	updateScoreRecods(records, &source[0], dataTelegram)

	err = s.storage.UpdateSource(&source[0])
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return err
	}

	return err
}

func (s *WafflerService) Search(search *message.Search) (*message.SearchResponse, error) {
	source, err := s.storage.SearchBySourceName(search.QueryForName)
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return nil, err
	}

	recordsRaw := make([]models.SourceDTO, 0, search.Limit)
	for i := search.Cursor; i < len(source) && i < search.Limit+search.Cursor; i++ {
		recordsRaw = append(recordsRaw, models.SourceDTO{
			Name:        source[i].Name,
			SourceType:  source[i].SourceType,
			SourceUrl:   source[i].SourceUrl,
			WaffelScore: source[i].WaffelScore,
			RacismScore: source[i].RacismScore,
		})
	}

	source = sortSources(recordsRaw, search.Order)
	search.Cursor += len(source)

	return &message.SearchResponse{
		Sources: source,
		Cursor:  search.Cursor,
	}, nil
}

func sortSources(sources []models.SourceDTO, search string) []models.SourceDTO {
	var orderSources = []string{"name_up", "name_down", "source_up", "source_down",
		"waffler_up", "waffler_down", "racism_up", "racism_down"}
	switch search {
	case orderSources[0]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].Name < sources[j].Name
		})
	case orderSources[1]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].Name > sources[j].Name
		})
	case orderSources[2]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].SourceType < sources[j].SourceType
		})
	case orderSources[3]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].SourceType > sources[j].SourceType
		})
	case orderSources[4]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].WaffelScore < sources[j].WaffelScore
		})
	case orderSources[5]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].WaffelScore > sources[j].WaffelScore
		})
	case orderSources[6]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].RacismScore < sources[j].RacismScore
		})
	case orderSources[7]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].RacismScore > sources[j].RacismScore
		})
	}

	return sources

}
