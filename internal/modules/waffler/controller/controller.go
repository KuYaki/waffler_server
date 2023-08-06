package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
)

type Waffler interface {
	Hello(http.ResponseWriter, *http.Request)
	HomePage(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
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
	//var data models.Search
	//var err error
	//
	//wa.Decoder.Decode(r.Body, &data)
	//if err != nil {
	//	wa.ErrorBadRequest(w, err)
	//	return
	//}

	//TODO implement me
	panic("implement me")
}

func (wa *Waffl) Hello(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("Hello, world!"))
	if err != nil {
		wa.log.Error("error writing response", zap.Error(err))
	}
}
