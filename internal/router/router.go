package router

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Recoverer) //  ToDO: Need?

	r.Get("/", controllers.Waffler.Hello)
	authCheck := components.Token

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
		r.Post("/search", sourceController.Search)
		r.Post("/parse", sourceController.Parse)
		r.HandleFunc("/ws", sourceController.WsTest)
		r.Post("/score", sourceController.Score)
		r.Post("/info", sourceController.Info)
	})

	return r
}
