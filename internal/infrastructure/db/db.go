package db

import (
	"database/sql"
	"fmt"

	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSqlDB(dbConf *config.AppConf) (*gorm.DB, error) {
	dsnRaw := "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable"

	dsn := fmt.Sprintf(dsnRaw,
		dbConf.DB.User, dbConf.DB.Password, dbConf.DB.Host, dbConf.DB.Port, dbConf.DB.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestDB(conn *gorm.DB) error {
	user := &models.UserDTO{}
	source := &models.SourceDTO{}
	records := &models.RecordDTO{}

	result := conn.Find(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		for i := 0; i < 100; i++ {
			userNew := &models.UserDTO{
				Username: gofakeit.Username(),
				PwdHash:  gofakeit.Animal(),
			}
			result = conn.Create(userNew)
			if result.Error != nil {
				return result.Error
			}

		}
	}

	result = conn.Find(source)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		for i := 0; i < 20; i++ {
			sourceNew := &models.SourceDTO{
				Name:         gofakeit.BeerName(),
				SourceType:   models.SourceType(gofakeit.Number(0, 1)),
				SourceUrl:    gofakeit.URL(),
				WafflerScore: models.NewNullFloat64(gofakeit.Float64Range(0, 100)),
				RacismScore:  models.NewNullFloat64(gofakeit.Float64Range(0, 100)),
			}
			result = conn.Create(sourceNew)
			if result.Error != nil {
				return result.Error
			}

		}
		for i := 0; i < 20; i++ {
			sourceNew := &models.SourceDTO{
				Name:         gofakeit.BeerName(),
				SourceType:   models.SourceType(gofakeit.Number(0, 1)),
				SourceUrl:    gofakeit.URL(),
				WafflerScore: models.NewNullFloat64(gofakeit.Float64Range(0, 100)),
			}
			result = conn.Create(sourceNew)
			if result.Error != nil {
				return result.Error
			}

		}
		for i := 0; i < 20; i++ {
			sourceNew := &models.SourceDTO{
				Name:        gofakeit.BeerName(),
				SourceType:  models.SourceType(gofakeit.Number(0, 1)),
				SourceUrl:   gofakeit.URL(),
				RacismScore: models.NewNullFloat64(gofakeit.Float64Range(0, 100)),
			}
			result = conn.Create(sourceNew)
			if result.Error != nil {
				return result.Error
			}

		}
	}

	racism := &models.RacismDTO{}
	result = conn.Find(racism)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		for i := 0; i < 10; i++ {
			racismNew := &models.RacismDTO{
				SourceID: gofakeit.Number(1, 3),
			}

			result = conn.Create(racismNew)
			if result.Error != nil {
				return result.Error
			}

		}
		for i := 0; i < 10; i++ {
			racismNew := &models.RacismDTO{
				SourceID: gofakeit.Number(1, 3),
				Score:    models.NewNullInt64(int64(gofakeit.Number(0, 100))),
			}

			result = conn.Create(racismNew)
			if result.Error != nil {
				return result.Error
			}

		}
	}

	result = conn.Find(records)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		for i := 0; i < 100; i++ {
			recordsNew := &models.RecordDTO{
				RecordText: gofakeit.Cat(),
				SourceID:   gofakeit.Number(1, 3),
			}

			result = conn.Create(recordsNew)
			if result.Error != nil {
				return result.Error
			}

		}
	}

	return nil
}
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
