package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"go.uber.org/zap"
)

type Waffler interface {
	Search(models.Search)
}

type WafflerService struct {
	storage storage.WafflerStorager
	log     *zap.Logger
}

func NewWafflerService(storage storage.WafflerStorager, components *component.Components) *WafflerService {
	return &WafflerService{storage: storage, log: components.Logger}
}

func (s *WafflerService) Search(search models.Search) {
	s.log.Info("search", zap.String("search", search.Query))
	//service.AnswerForGPT(s.log, search)
	return
}
