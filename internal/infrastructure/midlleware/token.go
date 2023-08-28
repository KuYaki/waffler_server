package midlleware

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/gin-gonic/gin"
	"log"

	"net/http"
	"strings"
)

const authorization = "Authorization"

type Token struct {
	responder.Responder
	jwt cryptography.TokenManager
}

type UserRequest struct{}

func NewTokenManager(responder responder.Responder, jwt cryptography.TokenManager) *Token {
	return &Token{
		Responder: responder,
		jwt:       jwt,
	}
}

func (t *Token) CheckStrictFunc(c *gin.Context) {
	tokenRaw := c.GetHeader(authorization)
	tokenParts := strings.Split(tokenRaw, " ")
	if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
		c.IndentedJSON(http.StatusForbidden, "wrong input data")
		return
	}
	_, err := t.jwt.ParseToken(tokenParts[1], cryptography.AccessToken)
	if err != nil && err.Error() == "Token is expired" {
		c.IndentedJSON(http.StatusUnauthorized, "token expired")
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (t *Token) CheckStrict() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRaw := c.GetHeader(authorization)
		log.Println(tokenRaw)
		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			c.IndentedJSON(http.StatusForbidden, "wrong input data")
			return
		}
		u, err := t.jwt.ParseToken(tokenParts[1], cryptography.AccessToken)
		if err != nil && err.Error() == "Token is expired" {
			c.IndentedJSON(http.StatusUnauthorized, "token expired")
			return
		}
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, err)
			return
		}

		ctx := context.WithValue(c.Request.Context(), UserRequest{}, u)
		c.Request.WithContext(ctx)
		c.Next()
	}
}

func (t *Token) CheckRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRaw := c.GetHeader(authorization)
		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			c.IndentedJSON(http.StatusForbidden, "wrong input data")
			return
		}
		u, err := t.jwt.ParseToken(tokenParts[1], cryptography.RefreshToken)
		if err != nil && err.Error() == "Token expired" {
			c.IndentedJSON(http.StatusUnauthorized, "token expired")
			return
		}
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, err)
			return
		}
		ctx := context.WithValue(context.Background(), UserRequest{}, u)
		c.Request.WithContext(ctx)
	}
}

func (t *Token) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(authorization)
		u, err := t.jwt.ParseToken(token, cryptography.AccessToken)
		if err != nil {
			u = cryptography.UserClaims{}
		}
		ctx := context.WithValue(context.Background(), UserRequest{}, u)
		c.Request.WithContext(ctx)
	}
}
