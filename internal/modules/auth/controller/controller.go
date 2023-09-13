package controller

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/handler"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/auth/service"
	"github.com/go-faster/errors"
	"github.com/ptflp/godecoder"

	"net/http"
)

type Auther interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type Auth struct {
	auth service.Auther
	responder.Responder
	godecoder.Decoder
}

func NewAuth(service service.Auther, components *component.Components) Auther {
	return &Auth{auth: service, Responder: components.Responder, Decoder: components.Decoder}
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := a.Decoder.Decode(r.Body, &req)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	if req.Password != req.RetypePassword {
		a.Responder.ErrorBadRequest(w, errors.New("passwords do not match"))
		return
	}

	statusCode, err := a.auth.Register(context.Background(), req.Username, req.Password)
	if err != nil {
		a.Error(w, statusCode, err)
		return
	}

	out, statusCode, err := a.auth.Login(context.Background(), models.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		a.Responder.Error(w, statusCode, err)
		return
	}

	a.Responder.OutputJSON(w, LoginData{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := a.Decoder.Decode(r.Body, &req)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}
	if req.Username == "" || req.Password == "" {
		a.Responder.ErrorBadRequest(w, errors.New("username and password are required"))
		return
	}

	out, statusCode, err := a.auth.Login(context.Background(), models.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		a.Responder.Error(w, statusCode, err)
		return
	}

	a.Responder.OutputJSON(w, LoginData{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}

func (a *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := handler.ExtractUser(r)
	if err != nil {
		a.Responder.ErrorBadRequest(w, err)
		return
	}
	out, err := a.auth.AuthorizeRefresh(context.Background(), claims.ID)
	if err != nil {
		a.Responder.ErrorBadRequest(w, err)
		return
	}

	a.Responder.OutputJSON(w, LoginData{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})

}
