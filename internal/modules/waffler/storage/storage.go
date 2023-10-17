package storage

import (
	"strings"

	"github.com/KuYaki/waffler_server/internal/models"
	"gorm.io/gorm"
)

type WafflerStorager interface {
	CreateSource(*models.SourceDTO) error
	CreateRecords(records []models.RecordDTO) error
	CreateRacismRecords(racismRecords []models.RacismDTO) error
	CreateWafflerRecords(racismRecords []models.WafflerDTO) error
	SearchLikeBySourceURLNotName(search string, sourceType []models.SourceType, offset int, order []string, limit int) ([]models.SourceDTO, error)
	SearchLikeBySourceName(search string, sourceType []models.SourceType, offset int, order []string, limit int) ([]models.SourceDTO, error)
	SearchBySourceUrl(url string) (*models.SourceDTO, error)
	UpdateSource(*models.SourceDTO) error
	ListRecordsSourceID(idSource int) ([]models.RecordDTO, error)
	ListRacismRecords(racismRecords *models.RacismDTO) ([]models.RacismDTO, error)
	ListWafflerRecords(racismRecords *models.WafflerDTO) ([]models.WafflerDTO, error)

	ListRecordsSourceIDCursor(idSource int, offset int, limit int) ([]models.RecordDTO, error)
	ListRacismRecordsSourceIDCursor(idSource int, order []string, offset int, limit int) ([]models.RacismDTO, error)
}

func NewWafflerStorage(conn *gorm.DB) WafflerStorager {
	return &WafflerStorage{conn: conn}
}

type WafflerStorage struct {
	conn *gorm.DB
}

func (s WafflerStorage) ListRacismRecordsSourceIDCursor(idSource int, order []string, offset int, limit int) ([]models.RacismDTO, error) {
	var records []models.RacismDTO

	querySQL := "source_id = ?"

	var args = []interface{}{idSource}

	db := ApplyOrder(s.conn, order)
	err := db.Offset(offset).Limit(limit).
		Where(querySQL, args...).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s WafflerStorage) ListRacismRecords(racismRecords *models.RacismDTO) ([]models.RacismDTO, error) {
	var records []models.RacismDTO
	err := s.conn.Where(racismRecords).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
func (s WafflerStorage) ListWafflerRecords(racismRecords *models.WafflerDTO) ([]models.WafflerDTO, error) {
	var records []models.WafflerDTO
	err := s.conn.Where(racismRecords).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s WafflerStorage) CreateRacismRecords(racismRecords []models.RacismDTO) error {
	err := s.conn.Create(racismRecords).Error
	if err != nil {
		return err
	}

	return nil
}

func (s WafflerStorage) CreateWafflerRecords(wafflerRecords []models.WafflerDTO) error {
	err := s.conn.Create(wafflerRecords).Error
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

func (s WafflerStorage) CreateRecords(source []models.RecordDTO) error {
	err := s.conn.Create(source).Error
	if err != nil {
		return err
	}

	return nil
}
func (s WafflerStorage) SearchLikeBySourceName(search string, sourceType []models.SourceType, offset int, order []string, limit int) ([]models.SourceDTO, error) {
	var sources []models.SourceDTO
	var querySQL string
	var args = make([]interface{}, 0, len(sourceType)+1)
	args = append(args, strings.ToLower("%"+search+"%"))
	if len(sourceType) == 1 {
		querySQL = "name ILIKE ? AND source_type = ?"
		args = append(args, sourceType[0])
	} else if len(sourceType) == 2 {
		querySQL = "name ILIKE ? AND (source_type = ? OR source_type = ?)"
		args = append(args, sourceType[0], sourceType[1])
	}

	db := ApplyOrder(s.conn, order)
	err := db.Offset(offset).Limit(limit).
		Where(querySQL, args...).Find(&sources).Error
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (s WafflerStorage) SearchLikeBySourceURLNotName(search string, sourceType []models.SourceType, offset int, order []string, limit int) ([]models.SourceDTO, error) {
	var sources []models.SourceDTO
	var querySQL string
	var args = make([]interface{}, 0, len(sourceType)+2)
	searchSQL := strings.ToLower("%" + search + "%")
	args = append(args, searchSQL, searchSQL)
	if len(sourceType) == 1 {
		querySQL = "source_url ILIKE ? AND not name ILIKE ? AND source_type = ?"
		args = append(args, sourceType[0])
	} else if len(sourceType) == 2 {
		querySQL = "source_url ILIKE ? AND not name ILIKE ? AND (source_type = ? OR source_type = ?)"
		args = append(args, sourceType[0], sourceType[1])
	}

	db := ApplyOrder(s.conn, order)
	err := db.Offset(offset).Limit(limit).
		Where(querySQL, args...).
		Find(&sources).Error
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (s WafflerStorage) SearchBySourceUrl(url string) (*models.SourceDTO, error) {
	source := models.SourceDTO{}

	res := s.conn.Where(models.SourceDTO{SourceUrl: url}).Find(&source)
	if res.Error != nil {
		return nil, res.Error
	}

	return &source, nil
}

func (s WafflerStorage) ListRecordsSourceID(idSource int) ([]models.RecordDTO, error) {
	var records []models.RecordDTO
	err := s.conn.Where(&models.RecordDTO{SourceID: idSource}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s WafflerStorage) ListRecordsSourceIDCursor(idSource int, offset int, limit int) ([]models.RecordDTO, error) {
	var records []models.RecordDTO

	querySQL := "source_id = ?"

	var args = []interface{}{idSource}

	err := s.conn.Offset(offset).Limit(limit).
		Where(querySQL, args...).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
