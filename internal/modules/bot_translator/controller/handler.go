package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator/service"
	"github.com/ptflp/godecoder"
	"net/http"
)

type BotTranslatorInterface interface {
	SetWebhook(w http.ResponseWriter, r *http.Request)
	DeleteWebhook(w http.ResponseWriter, r *http.Request)
}

type Translator struct {
	serviceTr service.BotTranslatorServicer
	responder.Responder
	godecoder.Decoder
}

func NewTranslator(service service.BotTranslatorServicer, components *component.Components) BotTranslatorInterface {
	return &Translator{serviceTr: service, Responder: components.Responder, Decoder: components.Decoder}

}

func (wa *Translator) SetWebhook(w http.ResponseWriter, r *http.Request) {
	var data bot_translator.SetWebhookParams
	err := wa.Decoder.Decode(r.Body, &data)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}
	var host = r.Host
	if r.Proto == "HTTP/1.1" {
		host = "http://" + host
	} else {
		host = "https://" + host
	}

	err = wa.serviceTr.SetWebhook(data, host)
	if err != nil {
		wa.ErrorInternal(w, err)
		return
	}

	wa.Responder.OutputJSON(w, nil)
}

func (wa *Translator) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	var data bot_translator.DeleteWebhookParams
	err := wa.Decoder.Decode(r.Body, &data)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	err = wa.serviceTr.DeleteWebhook(data, r.Host)
	if err != nil {
		wa.ErrorInternal(w, err)
		return
	}

	wa.Responder.OutputJSON(w, nil)

}
