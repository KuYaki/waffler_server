package mocks

import (
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/db"
	midle "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/KuYaki/waffler_server/internal/router"
	"github.com/KuYaki/waffler_server/internal/storages"
	"github.com/KuYaki/waffler_server/mocks/internal_/infrastructure/service/telegram"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"testing"
	"time"
)

type PostgreSQLContainer struct {
	Cont     testcontainers.Container
	ContInfo ContInfo
}

type ContInfo struct {
	UserName   string
	Password   string
	Host       string
	MappedPort string
	DbName     string
}

func MockServer(t *testing.T) *chi.Mux {
	z, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	appConf := &config.AppConf{
		AppName: "waffler_server",
		TestApp: false,
		Logger: &config.Logger{
			Level:   "debug",
			LogPath: "./app.log",
		},
		Server: &config.Server{
			Port: "8080",
		},
		DB: &config.DB{
			Name:     "postgres",
			User:     "postgres",
			Password: "postgres",
			Host:     "localhost",
			Port:     "5432",
		},
		Telegram: &config.Telegram{
			AppID:   0,
			ApiHash: "",
			Phone:   "",
		},
		ChatGPT: &config.ChatGPT{
			Token: "",
		},
		Token: &config.Token{
			AccessTTL:     time.Duration(10) * time.Minute,
			RefreshTTL:    time.Duration(50) * time.Minute,
			AccessSecret:  "secret",
			RefreshSecret: "bigsecret",
		},
	}

	srv := newMockApp(appConf, z, t)

	return srv
}

func newMockApp(conf *config.AppConf, logger *zap.Logger, t *testing.T) *chi.Mux {
	conn, err := db.NewSqlDB(conf)
	if err != nil {
		logger.Fatal("app: db error", zap.Error(err))
		return nil
	}
	err = conn.AutoMigrate(
		&models.UserDTO{},
		&models.SourceDTO{},
		&models.RecordDTO{},
		&models.RacismDTO{},
		&models.WafflerDTO{},
	)

	if conf.TestApp {
		err = db.TestDB(conn)
		if err != nil {
			logger.Fatal("app: db test error", zap.Error(err))
			return nil
		}
	}

	if err != nil {
		logger.Fatal("app: db migration error", zap.Error(err))
		return nil
	}
	tg := telegram.NewClientSource(t)
	if err != nil {
		logger.Fatal("app: tg error", zap.Error(err))
		return nil
	}
	dataSource := data_source.NewDataTelegram(tg)

	storagesDB := storages.NewStorages(conn)

	// инициализация менеджера токенов
	tokenManager := cryptography.NewTokenJWT(conf.Token)
	// инициализация декодера
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	// инициализация менеджера ответов сервера
	responseManager := responder.NewResponder(decoder, logger)
	// инициализация генератора uuid
	uuID := cryptography.NewUUIDGenerator()
	// инициализация хешера
	hash := cryptography.NewHash(uuID)

	token := midle.NewTokenManager(responseManager, tokenManager)

	components := component.NewComponents(conf, tokenManager, token, responseManager, decoder,
		hash, dataSource, logger)
	services := modules.NewServices(storagesDB, components)
	controller := modules.NewControllers(services, components)
	// init router
	r := router.NewApiRouter(controller, components)

	return r

}
