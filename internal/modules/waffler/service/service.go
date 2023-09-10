package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	tg "github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"go.uber.org/zap"
)

type Waffler interface {
	Search(search *models.Search) (*message.SearchResponse, error)
	InfoSource(domain string) *message.InfoRequest
	Score(request *message.ScoreRequest) (*message.ScoreResponse, error)
	ParseSource(search *message.ParserRequest) error
}

type WafflerService struct {
	storage storage.WafflerStorager
	log     *zap.Logger
	tg      *tg.Telegram
	gpt     *gpt.ChatGPT
}

func NewWafflerService(storage storage.WafflerStorager, components *component.Components) *WafflerService {
	return &WafflerService{storage: storage, log: components.Logger,
		gpt: components.Gpt, tg: components.Tg}
}

func (u *WafflerService) Score(request *message.ScoreRequest) (*message.ScoreResponse, error) {
	scoreResponse := &message.ScoreResponse{
		Cursor: request.Cursor,
	}

	records, err := u.storage.SelectRecords(request.SourceId)
	if err != nil {
		return nil, err
	}
	scoreResponse.Records = records

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

	for i, r := range dataTelegram.Records {
		if r.RecordText == "" {
			continue
		}

		dataTelegram.Records[i].Score, err = s.gpt.ConstructQuestionGPT(r.RecordText, search.ScoreType)
		if err != nil {
			s.log.Error("error: question gpt", zap.Error(err))
		}
		newRecords = append(newRecords, dataTelegram.Records[i])
	}

	dataTelegram.Records = newRecords

	err = s.storage.CreateSource(&dataTelegram.Source)
	if err != nil {
		s.log.Error("error: create", zap.Error(err))
	}

	return err
}

func (s *WafflerService) Search(search *models.Search) (*message.SearchResponse, error) {
	source, err := s.storage.SearchBySourceName(search)
	if err != nil {
		s.log.Error("error: search", zap.Error(err))
		return nil, err
	}

	return &message.SearchResponse{
		Sources: source,
		Cursor:  search.Cursor,
	}, nil
}
