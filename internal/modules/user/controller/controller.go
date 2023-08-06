package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/modules/user/service"
	"github.com/ptflp/godecoder"
	"net/http"
)

type Userer interface {
	Userinfo(w http.ResponseWriter, r *http.Request)
	UserSave(w http.ResponseWriter, r *http.Request)
}
type User struct {
	service service.Userer
	responder.Responder
	godecoder.Decoder
}

func NewUserController(service service.Userer, components *component.Components) Userer {
	return &User{service: service, Responder: components.Responder, Decoder: components.Decoder}
}

func (a *User) Userinfo(w http.ResponseWriter, r *http.Request) {

}

func (a *User) UserSave(w http.ResponseWriter, r *http.Request) {

}
