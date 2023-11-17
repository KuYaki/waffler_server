package controller

import (
	"fmt"

	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	middleware "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/models"
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
	ParseWebsocket(w http.ResponseWriter, r *http.Request)
	Price(w http.ResponseWriter, r *http.Request)
}

type Waffl struct {
	service     service.WafflerServicer
	log         *zap.Logger
	token       *middleware.Token
	userService service2.Userer
	responder.Responder
	godecoder.Decoder
}

func (wa *Waffl) GetField() service.WafflerServicer {
	return wa.service
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
func (wa *Waffl) ParseWebsocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}
	defer conn.Close()
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
	// Read request data from connection
	var Parser *message.ParserRequest
	err := wa.Decoder.Decode(r.Body, &Parser)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	if Parser.Parser.Type != models.YakiModel_GPT3_5TURBO && Parser.Parser.Token == "" {
		wa.Responder.ErrorBadRequest(w, err)
		return
	}

	var processedRecords float64
	var maxRecords int

	if Parser.ScoreType == models.Racism {
		maxRecords = Parser.Limit
	}
	if Parser.ScoreType == models.Waffler {
		maxRecords = Parser.Limit * (Parser.Limit - 1) / 2
	}

	//// Channel will be written to each time a record is processed
	updateChan := make(chan bool, maxRecords)

	err = wa.service.ParseSource(Parser, updateChan)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	for i := 0; i < maxRecords; i++ {
		<-updateChan
		processedRecords++
		progress := fmt.Sprintf("%.2f", processedRecords/float64(maxRecords))
		fmt.Println(progress)
		//err = conn.WriteMessage(websocket.TextMessage, []byte(progress)) // Send progress status to client
		if err != nil {
			wa.Responder.ErrorInternal(w, err)
			return
		}
	}

	wa.OutputJSON(w, nil)
}
func (wa *Waffl) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, world!")
	wa.Responder.OutputJSON(w, "Hello, world!")

}

func (wa *Waffl) Price(w http.ResponseWriter, r *http.Request) {
	var priceRequest *message.PriceRequest
	err := wa.Decoder.Decode(r.Body, &priceRequest)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	err = ValidatePriceRequest(*priceRequest)
	if err != nil {
		wa.ErrorBadRequest(w, err)
		return
	}

	priceResponse, err := wa.service.PriceSource(priceRequest)
	if err != nil {
		wa.Responder.ErrorInternal(w, err)
		return
	}

	wa.OutputJSON(w, priceResponse)
}
