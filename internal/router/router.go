package router

import (
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/go-chi/chi/v5"
)

func NewApiRouter(controllers *modules.Controllers) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", controllers.Waffler.Hello)

	return r
}
