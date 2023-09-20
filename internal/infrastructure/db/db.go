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

func NewSqlDB(dbConf config.AppConf) (*gorm.DB, error) {
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
				Username:    gofakeit.Username(),
				Hash:        gofakeit.Animal(),
				ParserToken: NewNullString(gofakeit.UUID()),
				ParserType:  NewNullString("GPT"),
				Locale:      NewNullString(gofakeit.RandomString([]string{"en", "ru"})),
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
		for i := 0; i < 100; i++ {
			sourceNew := &models.SourceDTO{
				Name:        gofakeit.BeerName(),
				SourceType:  0,
				SourceUrl:   gofakeit.URL(),
				WaffelScore: gofakeit.Number(0, 10),
				RacismScore: gofakeit.Number(0, 10),
			}
			result = conn.Create(sourceNew)
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
				Score:      gofakeit.Number(0, 10),
				ScoreType:  models.ScoreType(gofakeit.Number(0, 1)),
				CreatedAt:  gofakeit.Date(),
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
