package modules

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	acontroller "github.com/KuYaki/waffler_server/internal/modules/auth/controller"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator/controller"
	ucontroller "github.com/KuYaki/waffler_server/internal/modules/user/controller"
	wacontroller "github.com/KuYaki/waffler_server/internal/modules/waffler/controller"
)

type Controllers struct {
	Waffler wacontroller.Waffler
	Auth    acontroller.Auther
	User    ucontroller.Userer
	Bot     controller.BotTranslatorInterface
}

func NewControllers(services *Services, components *component.Components) *Controllers {
	authController := acontroller.NewAuth(services.Auth, components)
	wafflController := wacontroller.NewWaffl(services.WafflService, services.User, components)
	userController := ucontroller.NewUserController(services.User, components)
	botTranslator := controller.NewTranslator(services.BotService, components)

	return &Controllers{
		Waffler: wafflController,
		Auth:    authController,
		User:    userController,
		Bot:     botTranslator,
	}
}
