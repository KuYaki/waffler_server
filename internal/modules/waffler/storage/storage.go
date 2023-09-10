package storage

import (
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"gorm.io/gorm"
)

type WafflerStorager interface {
	CreateSource(*models.SourceDTO) error
	SearchBySourceName(search *models.Search) ([]message.Source, error)
	SelectRecords(idSource int) ([]models.RecordDTO, error)
}

func NewWafflerStorage(conn *gorm.DB) WafflerStorager {
	return &WafflerStorage{conn: conn}
}

type WafflerStorage struct {
	conn *gorm.DB
}

func (s WafflerStorage) CreateSource(source *models.SourceDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) CreateRecord(source *models.SourceDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) SearchBySourceName(search *models.Search) ([]message.Source, error) {
	sources := []models.SourceDTO{}

	err := s.conn.Where("name LIKE ?", "%"+search.QueryForName+"%").Find(&sources).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s WafflerStorage) SelectRecords(idSource int) ([]models.RecordDTO, error) {
	records := []models.RecordDTO{}
	err := s.conn.Where(&models.RecordDTO{SourceID: idSource}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
