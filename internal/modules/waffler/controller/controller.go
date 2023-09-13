package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"

	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Waffler interface {
	Hello(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Score(w http.ResponseWriter, r *http.Request)
	Info(w http.ResponseWriter, r *http.Request)
	Parse(w http.ResponseWriter, r *http.Request)
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

func (wa *Waffl) Search(w http.ResponseWriter, r *http.Request) {
	var data models.Search
	err := wa.Decoder.Decode(r.Body, &data)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}
	search, err := wa.service.Search(&data)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	wa.Responder.OutputJSON(w, search)
}

func (wa *Waffl) Score(w http.ResponseWriter, r *http.Request) {
	var scoreRequest message.ScoreRequest
	err := wa.Decoder.Decode(r.Body, &scoreRequest)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	scoreResponse, err := wa.service.Score(&scoreRequest)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	wa.Responder.OutputJSON(w, scoreResponse)
}
func (wa *Waffl) Info(w http.ResponseWriter, r *http.Request) {
	var sourceURL message.SourceURL
	err := wa.Decoder.Decode(r.Body, &sourceURL)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	u, err := url.Parse(sourceURL.SourceUrl)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	info := wa.service.InfoSource(u.Hostname())
	if info == nil {
		wa.Responder.ErrorBadRequest(w, err)
		return
	}

	wa.Responder.OutputJSON(w, info)
}

func (wa *Waffl) Parse(w http.ResponseWriter, r *http.Request) {
	var Parser *message.ParserRequest
	err := wa.Decoder.Decode(r.Body, &Parser)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	err = wa.service.ParseSource(Parser)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	wa.OutputJSON(w, nil)
}
func (wa *Waffl) Hello(w http.ResponseWriter, r *http.Request) {
	wa.Responder.OutputJSON(w, "Hello, world!")

}
