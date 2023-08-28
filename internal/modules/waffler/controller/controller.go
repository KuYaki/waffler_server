package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/gin-gonic/gin"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Waffler interface {
	Hello(c *gin.Context)
	HomePage(c *gin.Context)
	Search(c *gin.Context)
	Score(c *gin.Context)
	Info(c *gin.Context)
	Parse(c *gin.Context)
}

type Waffl struct {
	service service.Waffler
	log     *zap.Logger
	responder.Responder
	godecoder.Decoder
}

func NewWaffl(service service.Waffler, components *component.Components) Waffler {
	return &Waffl{service: service,
		log: components.Logger, Responder: components.Responder, Decoder: components.Decoder}
}

func (wa *Waffl) HomePage(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (wa *Waffl) Search(c *gin.Context) {
	data := &models.Search{}
	err := c.BindJSON(data)
	if err != nil {
		return
	}

	_, err = url.ParseRequestURI(data.Query)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	wa.log.Info("search", zap.String("search", data.Query))

}

func (wa *Waffl) Score(c *gin.Context) {
	data := &models.Source{}
	c.BindJSON(data)

}
func (wa *Waffl) Info(c *gin.Context) {

}

func (wa *Waffl) Parse(c *gin.Context) {

	//wa.service.Search(data)
}
func (wa *Waffl) Hello(c *gin.Context) {
	_, err := c.Writer.Write([]byte("Hello, world!"))
	if err != nil {
		wa.log.Error("error writing response", zap.Error(err))
	}
}
