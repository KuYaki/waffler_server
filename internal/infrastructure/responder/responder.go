package responder

import (
	"context"
	"errors"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easytags $GOFILE
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{}, statusCode int)

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type Respond struct {
	log *zap.Logger
	godecoder.Decoder
}

func NewResponder(decoder godecoder.Decoder, logger *zap.Logger) Responder {
	return &Respond{log: logger, Decoder: decoder}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	if err := r.Encode(w, responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http response forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http   response Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
}
