package run

import (
	"context"
	"fmt"
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/db"
	midle "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/server"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/KuYaki/waffler_server/internal/router"
	"github.com/KuYaki/waffler_server/internal/storages"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

const (
	NoError = iota
	GeneralError
)

// Application - application interface
type Application interface {
	Runner
	Bootstraper
}

// Runner - application launch interface
type Runner interface {
	Run() int
}

// Bootstraper - application initialization interface
type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

// App - application structure
type App struct {
	conf   config.AppConf
	logger *zap.Logger
	srv    server.Server
	Sig    chan os.Signal
}

// NewApp - application builder
func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

// Run - application launch
func (a *App) Run() int {
	// create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	// start goroutines for graceful shutdown
	// when a signal is received SIGINT
	// call cancel for context
	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt received", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return GeneralError
	}

	return NoError
}

// Bootstrap - init application
func (a *App) Bootstrap() Runner {
	conn, err := db.NewSqlDB(a.conf)
	if err != nil {
		a.logger.Fatal("app: db error", zap.Error(err))
		return nil
	}
	err = conn.AutoMigrate(
		&models.UserDTO{},
		&models.WafflerDTO{},
		&models.WafflerDTO{},
		&models.RacismDTO{},
	)

	if a.conf.TestApp {
		err = db.TestDB(conn)
		if err != nil {
			a.logger.Fatal("app: db test error", zap.Error(err))
			return nil
		}
	}

	if err != nil {
		a.logger.Fatal("app: db migration error", zap.Error(err))
		return nil
	}
	tg, err := telegram.NewTelegram(a.conf.Telegram)
	if err != nil {
		a.logger.Fatal("app: tg error", zap.Error(err))
		return nil
	}
	dataSource := data_source.NewDataTelegram(tg)

	storagesDB := storages.NewStorages(conn)

	// инициализация менеджера токенов
	tokenManager := cryptography.NewTokenJWT(a.conf.Token)
	// инициализация декодера
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	// инициализация менеджера ответов сервера
	responseManager := responder.NewResponder(decoder, a.logger)
	// инициализация генератора uuid
	uuID := cryptography.NewUUIDGenerator()
	// инициализация хешера
	hash := cryptography.NewHash(uuID)

	token := midle.NewTokenManager(responseManager, tokenManager)

	components := component.NewComponents(a.conf, tokenManager, token, responseManager, decoder,
		hash, dataSource, a.logger)
	services := modules.NewServices(storagesDB, components)
	controller := modules.NewControllers(services, components)
	// init router
	r := router.NewApiRouter(controller, components)
	// server configuration
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}
	// server initialization
	a.srv = server.NewHttpServer(a.conf.Server, srv, a.logger)

	return a
}
