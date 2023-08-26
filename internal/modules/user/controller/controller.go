package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/modules/user/service"
	"github.com/gin-gonic/gin"
	"github.com/ptflp/godecoder"
)

type Userer interface {
	Info(c *gin.Context)
	Save(c *gin.Context)
}
type User struct {
	service service.Userer
	responder.Responder
	godecoder.Decoder
}

func NewUserController(service service.Userer, components *component.Components) Userer {
	return &User{service: service, Responder: components.Responder, Decoder: components.Decoder}
}

func (a *User) Info(c *gin.Context) {

}

func (a *User) Save(c *gin.Context) {

}
