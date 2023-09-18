package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	tg "github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/zap"
	"sort"
	"time"
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
	scoreResponse := &message.ScoreResponse{
		Cursor: request.Cursor,
	}

	records, err := u.storage.SelectRecords(request.SourceId)
	if err != nil {
		return nil, err
	}
	for i := scoreResponse.Cursor; i < len(records) && i < request.Limit+scoreResponse.Cursor+request.Limit; i++ {

		scoreResponse.Records = append(scoreResponse.Records, message.Record{
			RecordText: records[i].RecordText,
			Score:      records[i].Score,
			Timestamp:  records[i].CreatedAt,
		})
	}

	goFake := make([]message.Record, 0, 101)

	for i := 0; i < 100; i++ {
		s := message.Record{
			RecordText: gofakeit.Cat(),
			Score:      gofakeit.IntRange(0, 10),
			Timestamp:  time.Now(),
		}
		goFake = append(goFake, s)
	}
	scoreResponse.Records = goFake

	scoreResponse.Records = scoreResponse.Records[scoreResponse.Cursor : scoreResponse.Cursor+request.Limit]
	scoreResponse.Cursor += len(scoreResponse.Records)

	return scoreResponse, nil
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
	goFake := make([]models.SourceDTO, 0, 101)

	for i := 0; i < 100; i++ {
		s := models.SourceDTO{
			ID:          i,
			Name:        gofakeit.Name(),
			SourceType:  models.SourceType(gofakeit.RandomInt([]int{0, 1})),
			SourceUrl:   gofakeit.URL(),
			WaffelScore: gofakeit.Number(0, 10),
			RacismScore: gofakeit.Number(0, 10),
		}
		goFake = append(goFake, s)
	}
	source = sortSources(goFake, search)
	sourceRes := source[search.Cursor : search.Cursor+search.Limit]
	search.Cursor += len(sourceRes)

	return &message.SearchResponse{
		Sources: sourceRes,
		Cursor:  search.Cursor,
	}, nil
}

func sortSources(sources []models.SourceDTO, search *message.Search) []models.SourceDTO {
	switch search.Order {
	case orderSort[0]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].Name > sources[j].Name
		})
	case orderSort[1]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].Name < sources[j].Name
		})
	case orderSort[2]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].SourceType > sources[j].SourceType
		})
	case orderSort[3]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].SourceType < sources[j].SourceType
		})
	case orderSort[4]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].WaffelScore > sources[j].WaffelScore
		})
	case orderSort[5]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].WaffelScore < sources[j].WaffelScore
		})
	case orderSort[6]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].RacismScore > sources[j].RacismScore
		})
	case orderSort[7]:
		sort.Slice(sources, func(i, j int) bool {
			return sources[i].RacismScore < sources[j].RacismScore
		})
	}

	return sources

}

var orderSort = []string{"name_up", "name_down", "source_up", "source_down",
	"waffler_up", "waffler_down", "racism_up", "racism_down"}
