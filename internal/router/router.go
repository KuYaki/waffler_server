package router

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/go-chi/chi/v5"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) *chi.Mux {
	r := chi.NewRouter()
	authCheck := midlleware.NewTokenManager(components.Responder, components.TokenManager)

	r.Get("/", controllers.Waffler.Hello)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.Auth.Register)
		r.Post("/login", controllers.Auth.Login)
		r.Route("/refresh", func(r chi.Router) {
			r.Use(authCheck.CheckRefresh)
			r.Post("/", controllers.Auth.Refresh)
		})

	})

	r.Route("/user", func(r chi.Router) {
		userController := controllers.User
		r.Post("/save", userController.Save)
		r.Post("/info", userController.Info)
		//r.Route("/profile", func(r chi.Router) {
		//	r.Use(authCheck.CheckStrict)
		//	r.Get("/", userController.Profile)
		//	r.Post("/changed_password", userController.ChangePassword)
		//})
	})

	r.Route("/source", func(r chi.Router) {
		r.Use(authCheck.CheckStrict)
		r.Post("/search", controllers.Waffler.Search)
		r.Post("/score", controllers.Waffler.Score)
		r.Post("info", controllers.Waffler.Info)
		r.Post("/parse", controllers.Waffler.Parse)
	})

	return r
}
