package responder

import (
	"context"
	"errors"
	"github.com/ptflp/godecoder"
	"net/http"

	"go.uber.org/zap"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
	ErrorConflict(w http.ResponseWriter, err error)
	Error(w http.ResponseWriter, status int, err error)
}

type Respond struct {
	log *zap.Logger
	godecoder.Decoder
}

func NewResponder(decoder godecoder.Decoder, logger *zap.Logger) Responder {
	return &Respond{log: logger, Decoder: decoder}
}

func (r *Respond) Error(w http.ResponseWriter, status int, err error) {
	switch status {
	case http.StatusBadRequest:
		r.ErrorBadRequest(w, err)
	case http.StatusUnauthorized:
		r.ErrorUnauthorized(w, err)
	case http.StatusForbidden:
		r.ErrorForbidden(w, err)
	case http.StatusConflict:
		r.ErrorConflict(w, err)
	case http.StatusInternalServerError:
		r.ErrorInternal(w, err)

	default:
		r.log.Error("responder error: unknown status code", zap.Error(err))
	}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if responseData != nil {
		err := r.Encode(w, responseData)
		if err != nil {
			r.log.Error("responder json encode error", zap.Error(err))
		}
	}
}

func (r *Respond) ErrorConflict(w http.ResponseWriter, err error) {
	r.log.Warn("http response conflict", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusConflict)
	if err := r.Encode(w, err.Error()); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := r.Encode(w, err.Error()); err != nil {
		r.log.Info("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http response forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	if err := r.Encode(w, err.Error()); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http   response Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := r.Encode(w, err.Error()); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := r.Encode(w, err.Error()); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}
