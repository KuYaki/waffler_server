package modules

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	aservice "github.com/KuYaki/waffler_server/internal/modules/auth/service"
	uservice "github.com/KuYaki/waffler_server/internal/modules/user/service"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/KuYaki/waffler_server/internal/storages"
)

type Services struct {
	User         uservice.Userer
	Auth         aservice.Auther
	WafflService service.WafflerServicer
}

func NewServices(storages *storages.Storages, components *component.Components) *Services {
	userService := uservice.NewUserService(storages.User, components)
	return &Services{
		User:         userService,
		Auth:         aservice.NewAuthService(userService, components),
		WafflService: service.NewWafflerService(storages.Waffler, components)}
}
