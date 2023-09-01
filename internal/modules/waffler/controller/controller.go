package controller

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/gin-gonic/gin"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Waffler interface {
	Hello(c *gin.Context)
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

func (wa *Waffl) Search(c *gin.Context) {
	var data models.Search
	err := c.BindJSON(&data)
	if err != nil {
		return
	}
	search, err := wa.service.Search(&data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, search)
}

func (wa *Waffl) Score(c *gin.Context) {
	var scoreRequest message.ScoreRequest
	err := c.BindJSON(&scoreRequest)
	if err != nil {
		return
	}

	scoreResponse, err := wa.service.Score(&scoreRequest)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, scoreResponse)
}
func (wa *Waffl) Info(c *gin.Context) {
	var sourceURL message.SourceURL
	err := c.BindJSON(&sourceURL)
	if err != nil {
		return
	}

	u, err := url.Parse(sourceURL.SourceUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	info := wa.service.InfoSource(u.Hostname())
	if info == nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	c.IndentedJSON(http.StatusOK, info)
}

func (wa *Waffl) Parse(c *gin.Context) {
	var Parser *message.ParserRequest
	err := c.BindJSON(&Parser)
	if err != nil {
		return
	}

	err = wa.service.ParseSource(Parser)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
func (wa *Waffl) Hello(c *gin.Context) {
	_, err := c.Writer.Write([]byte("Hello, world!"))
	if err != nil {
		wa.log.Error("error writing response", zap.Error(err))
	}
}
