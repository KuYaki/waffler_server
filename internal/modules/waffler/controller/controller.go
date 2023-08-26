package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Waffler interface {
	Hello(http.ResponseWriter, *http.Request)
	HomePage(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Score(http.ResponseWriter, *http.Request)
	Info(http.ResponseWriter, *http.Request)
	Parse(http.ResponseWriter, *http.Request)
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

func (wa *Waffl) HomePage(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (wa *Waffl) Search(w http.ResponseWriter, r *http.Request) {
	data := &models.Search{}
	err := wa.Decoder.Decode(r.Body, data)
	if err != nil {
		wa.Responder.ErrorBadRequest(w, err)
		return
	}

	_, err = url.ParseRequestURI(data.Query)
	if err != nil {
		wa.Responder.ErrorBadRequest(w, err)
		return
	}

	wa.log.Info("search", zap.String("search", data.Query))

}

func (wa *Waffl) Score(w http.ResponseWriter, r *http.Request) {
	data := &models.Score{}
	err := wa.Decoder.Decode(r.Body, data)
	if err != nil {
		wa.Responder.ErrorBadRequest(w, err)
		return
	}

}
func (wa *Waffl) Info(w http.ResponseWriter, r *http.Request) {

}

func (wa *Waffl) Parse(w http.ResponseWriter, r *http.Request) {

	//wa.service.Search(data)
}
func (wa *Waffl) Hello(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("Hello, world!"))
	if err != nil {
		wa.log.Error("error writing response", zap.Error(err))
	}
}
