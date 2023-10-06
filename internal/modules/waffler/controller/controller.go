package controller

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	middleware "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	service2 "github.com/KuYaki/waffler_server/internal/modules/user/service"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/gorilla/websocket"

	"net/http"

	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Waffler interface {
	Hello(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Score(w http.ResponseWriter, r *http.Request)
	Info(w http.ResponseWriter, r *http.Request)
	Parse(w http.ResponseWriter, r *http.Request)
	WsTest(w http.ResponseWriter, r *http.Request)
}

type Waffl struct {
	service     service.WafflerServicer
	log         *zap.Logger
	token       *middleware.Token
	userService service2.Userer
	responder.Responder
	godecoder.Decoder
}

func NewWaffl(service service.WafflerServicer, user service2.Userer, components *component.Components) Waffler {
	return &Waffl{service: service,
		log: components.Logger, token: components.Token, userService: user, Responder: components.Responder, Decoder: components.Decoder}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (wa *Waffl) WsTest(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	err = reader(ws)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

}

func reader(conn *websocket.Conn) error {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			return err
		}
	}

}

func (wa *Waffl) Search(w http.ResponseWriter, r *http.Request) {
	var data message.Search
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

	info, err := wa.service.InfoSource(sourceURL.SourceUrl)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
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

	if Parser.Parser.Token == "" || Parser.Parser.Type == "" {
		wa.Responder.ErrorBadRequest(w, err)
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
