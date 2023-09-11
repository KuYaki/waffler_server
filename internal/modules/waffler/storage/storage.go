package storage

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"gorm.io/gorm"
)

type WafflerStorager interface {
	CreateSource(*models.SourceDTO) error
	SearchBySourceName(search *models.Search) ([]models.SourceDTO, error)
	CreateSourceAndRecords(source *telegram.DataTelegram) error
	SelectRecords(idSource int) ([]models.RecordDTO, error)
}

func NewWafflerStorage(conn *gorm.DB) WafflerStorager {
	return &WafflerStorage{conn: conn}
}

type WafflerStorage struct {
	conn *gorm.DB
}

func (s WafflerStorage) CreateSourceAndRecords(source *telegram.DataTelegram) error {
	tx := s.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Create(&source.Source).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range source.Records {
		source.Records[i].SourceID = source.Source.ID
	}

	err = tx.Create(source.Records).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) CreateSource(source *models.SourceDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) CreateRecords(source []models.RecordDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) SearchBySourceName(search *models.Search) ([]models.SourceDTO, error) {
	sources := []models.SourceDTO{}

	err := s.conn.Where("name LIKE ?", "%"+search.QueryForName+"%").Find(&sources).Error
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (s WafflerStorage) SelectRecords(idSource int) ([]models.RecordDTO, error) {
	records := []models.RecordDTO{}
	err := s.conn.Where(&models.RecordDTO{SourceID: idSource}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
