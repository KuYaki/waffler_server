package modules

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	acontroller "github.com/KuYaki/waffler_server/internal/modules/auth/controller"
	ucontroller "github.com/KuYaki/waffler_server/internal/modules/user/controller"
	wacontroller "github.com/KuYaki/waffler_server/internal/modules/waffler/controller"
)

type Controllers struct {
	Waffler wacontroller.Waffler
	Auth    acontroller.Auther
	User    ucontroller.Userer
}

func NewControllers(services *Services, components *component.Components) *Controllers {
	authController := acontroller.NewAuth(services.Auth, components)
	wafflController := wacontroller.NewWaffl(services.WafflService, components)
	userController := ucontroller.NewUserController(services.User, components)

	return &Controllers{
		Waffler: wafflController,
		Auth:    authController,
		User:    userController,
	}
}
