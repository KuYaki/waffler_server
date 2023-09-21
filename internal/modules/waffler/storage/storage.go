package storage

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"gorm.io/gorm"
	"strings"
)

type WafflerStorager interface {
	CreateSource(*models.SourceDTO) error

	CreateRecords(records []*models.RecordDTO) error
	SearchByLikeSourceName(search string, cursor message.Cursor, order string, limit int) ([]models.SourceDTO, *message.Cursor, error)
	SearchBySourceUrl(url string) (*models.SourceDTO, error)
	UpdateSource(*models.SourceDTO) error
	CreateSourceAndRecords(source *telegram.DataTelegram) error
	SelectRecordsSourceID(idSource int) ([]*models.RecordDTO, error)
	SelectRecordsSourceIDOffsetLimit(idSource int, offset int, limit int) ([]models.RecordDTO, error)
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

func (s WafflerStorage) UpdateSource(source *models.SourceDTO) error {
	err := s.conn.Model(source).Updates(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) CreateRecords(source []*models.RecordDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) SearchByLikeSourceName(search string, cursor message.Cursor, order string, limit int) ([]models.SourceDTO, *message.Cursor, error) {
	var sources, sourcesURL, sourcesName []models.SourceDTO
	var err error
	searchSQL := strings.ToLower("%" + search + "%")

	if cursor.Partition == 0 {
		sourcesName, cursor, limit, err = s.searchLikeSources("name ILIKE ?", cursor, limit, order, searchSQL)
		if err != nil {
			return nil, nil, err
		}
		sources = append(sources, sourcesName...)

	}

	if cursor.Partition == 1 {
		sourcesURL, cursor, _, err = s.searchLikeSources("source_url ILIKE ? and not name ILIKE ?", cursor, limit, order, searchSQL, searchSQL)
		if err != nil {
			return nil, nil, err
		}
		cursor.Partition = 1
		sources = append(sources, sourcesURL...)
	}

	return sources, &cursor, nil
}

func (s WafflerStorage) searchLikeSources(query string, cursor message.Cursor, limit int, order string, arg ...interface{}) ([]models.SourceDTO, message.Cursor, int, error) {
	var records []models.SourceDTO
	resURL := s.conn.Order(order).Offset(cursor.Offset).Limit(limit).Where(query, arg...).Find(&records)
	if resURL.Error != nil {
		return nil, message.Cursor{}, 0, resURL.Error
	}

	cursor.Offset += int(resURL.RowsAffected)
	limit -= int(resURL.RowsAffected)
	if limit > 0 {
		cursor.Offset = 0
		cursor.Partition += 1
	}

	return records, cursor, limit, nil
}

func (s WafflerStorage) SearchBySourceUrl(url string) (*models.SourceDTO, error) {
	source := models.SourceDTO{}

	res := s.conn.Where(models.SourceDTO{SourceUrl: url}).Find(&source)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil

	}

	return &source, nil
}

func (s WafflerStorage) SelectRecordsSourceID(idSource int) ([]*models.RecordDTO, error) {
	var records []*models.RecordDTO
	err := s.conn.Where(&models.RecordDTO{SourceID: idSource}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s WafflerStorage) SelectRecordsSourceIDOffsetLimit(idSource int, offset int, limit int) ([]models.RecordDTO, error) {
	var records []models.RecordDTO
	err := s.conn.Offset(offset).Limit(limit).Where(&models.RecordDTO{SourceID: idSource}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
