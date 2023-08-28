package router

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/modules"
	"github.com/gin-gonic/gin"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) *gin.Engine {
	newR := gin.Default()
	authCheck := midlleware.NewTokenManager(components.Responder, components.TokenManager)

	newR.GET("/", controllers.Waffler.Hello)

	auth := newR.Group("/auth")
	{
		auth.POST("/register", controllers.Auth.Register)
		auth.POST("/login", controllers.Auth.Login)
		auth.Group("/refresh")
		{
			auth.Use(authCheck.CheckRefresh())
			auth.POST("/", controllers.Auth.Refresh)

		}

	}

	user := newR.Group("/user")
	{
		userController := controllers.User
		user.Use(authCheck.CheckStrict())
		user.POST("/save", userController.Save)
		user.POST("/info", userController.Info)
		//r.Route("/profile", func(r chi.Router) {
		//	r.Use(authCheck.CheckStrict)
		//	r.Get("/", userController.Profile)
		//	r.Post("/changed_password", userController.ChangePassword)
		//})
	}

	source := newR.Group("/source")
	{
		sourceController := controllers.Waffler
		source.Use(authCheck.CheckStrict())
		source.POST("/search", sourceController.Search)
		source.POST("/score", sourceController.Score)
		source.POST("/info", sourceController.Info)
		source.POST("/parse", sourceController.Parse)
	}

	return newR
}
