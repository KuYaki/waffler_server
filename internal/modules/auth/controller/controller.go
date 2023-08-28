package controller

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/handler"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/ptflp/godecoder"

	"net/http"
)

type Auther interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type Auth struct {
	auth service.Auther
	responder.Responder
	godecoder.Decoder
}

func NewAuth(service service.Auther, components *component.Components) Auther {
	return &Auth{auth: service, Responder: components.Responder, Decoder: components.Decoder}
}

func (a *Auth) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if req.Password != req.RetypePassword {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "passwords don't match",
		})

		return
	}

	err := a.auth.Register(context.Background(), req.Username, req.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (a *Auth) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}
	if req.Username == "" || req.Password == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	out, err := a.auth.Login(context.Background(), models.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, LoginData{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}

func (a *Auth) Refresh(c *gin.Context) {
	claims, err := handler.ExtractUser(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	out, err := a.auth.AuthorizeRefresh(context.Background(), claims.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, LoginData{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})

}
