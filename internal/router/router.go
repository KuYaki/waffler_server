package router

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	midle "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", controllers.Waffler.Hello)
	authCheck := midle.NewTokenManager(components.Responder, components.TokenManager)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.Auth.Register)
		r.Post("/login", controllers.Auth.Login)

		r.Route("/refresh", func(r chi.Router) {
			r.Use(authCheck.CheckRefresh)
			r.Post("/", controllers.Auth.Refresh)
		})

	})
	r.Route("/user", func(r chi.Router) {
		r.Use(authCheck.CheckStrict)
		r.Get("/info", controllers.User.Info)
		r.Post("/save", controllers.User.Save)
	})

	r.Route("/source", func(r chi.Router) {
		sourceController := controllers.Waffler
		r.Use(authCheck.CheckStrict)
		r.Post("/search", sourceController.Search)
		r.Post("/score", sourceController.Score)
		r.Post("/info", sourceController.Info)
		r.Post("/parse", sourceController.Parse)
	})

	return r
}
