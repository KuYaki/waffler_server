package controller

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/errors"
	"github.com/KuYaki/waffler_server/internal/infrastructure/handler"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/auth/service"
	"github.com/ptflp/godecoder"

	"net/http"
)

type Auther interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
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
	err := a.Decode(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Password != req.RetypePassword {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errorCode, myErr := a.auth.Register(context.Background(), req.Username, req.Password)
	if myErr != errors.NoError {
		w.WriteHeader(errorCode)
		return
	}

	w.WriteHeader(errorCode)
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := a.Decode(r.Body, &req)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out := a.auth.Login(r.Context(), models.User{
		Username: req.Username,
		Hash:     req.Password,
	})

	if out.ErrorCode != errors.NoError {
		a.OutputJSON(w, "", http.StatusUnauthorized)

		return
	}

	a.OutputJSON(w,
		LoginData{
			AccessToken:  out.AccessToken,
			RefreshToken: out.RefreshToken,
		}, http.StatusOK)
}

func (a *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := handler.ExtractUser(r)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}
	out := a.auth.AuthorizeRefresh(r.Context(), claims.ID)

	if out.ErrorCode != errors.NoError {
		a.OutputJSON(w, "", http.StatusUnauthorized)
		return
	}

	a.OutputJSON(w,
		LoginData{
			AccessToken:  out.AccessToken,
			RefreshToken: out.RefreshToken,
		}, http.StatusOK)
}
