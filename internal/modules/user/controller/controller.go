package controller

import (
	"context"
	"errors"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/handler"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/user/service"
	"github.com/ptflp/godecoder"
	"net/http"
)

type Userer interface {
	Info(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
}
type User struct {
	service service.Userer
	jwt     cryptography.TokenManager
	responder.Responder
	godecoder.Decoder
}

func NewUserController(service service.Userer, components *component.Components) Userer {
	return &User{service: service, Responder: components.Responder, Decoder: components.Decoder, jwt: components.TokenManager}
}

func (a *User) Info(w http.ResponseWriter, r *http.Request) {
	claims, err := handler.ExtractUser(r)
	if err != nil {
		a.ErrorInternal(w, err)
		return
	}
	userInfo, err := a.service.GetUserInfo(context.Background(), claims.ID)
	if err != nil {
		a.ErrorInternal(w, err)
		return
	}

	a.Responder.OutputJSON(w, userInfo)
}

func validate(user message.UserInfo) bool {

	return message.ValidateLocale(user.Locale) && message.ValidateParser(int(user.Parser.Type))

}

func (a *User) Save(w http.ResponseWriter, r *http.Request) {
	claims, err := handler.ExtractUser(r)
	if err != nil {
		a.ErrorInternal(w, err)
		return
	}
	var userSave message.UserInfo
	err = a.Decoder.Decode(r.Body, &userSave)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	if !validate(userSave) {
		a.Responder.ErrorBadRequest(w, errors.New("invalid locale or parser"))
		return
	}

	err = a.service.Update(context.Background(), userSave, claims.ID)
	if err != nil {
		a.Responder.ErrorInternal(w, err)
		return
	}

	a.OutputJSON(w, nil)
}
