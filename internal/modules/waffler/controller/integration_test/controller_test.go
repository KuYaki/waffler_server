//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/language_model"
	"github.com/KuYaki/waffler_server/mocks"
	"github.com/goccy/go-json"
	"github.com/gotd/td/tg"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestWaffl_Info(t *testing.T) {
	type mocksArgs struct {
		arguments string
		channel   []tg.ChatClass
		error     error
		isCall    bool
	}
	type result struct {
		name       string
		typeSource models.SourceType
		code       int
	}

	tests := []struct {
		name      string
		request   string
		mocksArgs mocksArgs
		result    result
	}{
		// TODO: Add test cases.
		{
			name:    "default telegram",
			request: "https://t.me/maximkatz",
			mocksArgs: mocksArgs{
				arguments: "https://t.me/maximkatz",
				channel:   []tg.ChatClass{&tg.Channel{Username: "maximkatz", Title: "Канал Максима Каца"}},
				error:     nil,
				isCall:    true,
			},
			result: result{name: "Канал Максима Каца",
				typeSource: models.Telegram,
				code:       http.StatusOK},
		},
		{
			name:    "unknown",
			request: "",
			mocksArgs: mocksArgs{
				isCall: false,
			},
			result: result{name: "Unknown",
				typeSource: models.Unknown,
				code:       http.StatusOK},
		},
	}

	srv, srvComponents := mocks.NewMockServer(t)
	rr := httptest.NewRecorder()
	assertTest := assert.New(t)
	i := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mocksArgs.isCall {
				srvComponents.MockComponents.TgClient.On("ContactSearch", tt.mocksArgs.arguments).
					Return(&tg.ContactsFound{Chats: tt.mocksArgs.channel}, nil).Once()
			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(message.SourceURL{SourceUrl: tt.request})
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/source/info", &buf)
			if err != nil {
				t.Fatal(err)
			}

			// Использование httptest для записи ответа
			srv.ServeHTTP(rr, req)

			res := message.InfoRequest{}
			err = json.Unmarshal(bytes.Split(rr.Body.Bytes(), []byte("\n"))[i], &res)
			if err != nil {
				t.Fatal(err)
			}
			assertTest.Equal(tt.result.code, rr.Code)

			assertTest.Equal(tt.result.name, res.Name)
			assertTest.Equal(tt.result.typeSource, res.Type)

			i++
		})
	}
}

func TestWaffl_Parse(t *testing.T) {
	type argumentsTg struct {
		methodName string
		args       []interface{}
		argsReturn []interface{}
		isCall     bool
	}
	type argumentsGpt struct {
		methodName string
		args       []interface{}
		argsReturn []interface{}
		isCall     bool
	}
	type mocksArgs struct {
		argumentsTg  []argumentsTg
		argumentsGpt []argumentsGpt
	}
	type result struct {
		name       string
		typeSource models.SourceType
		code       int
	}

	messages := func(amount int) tg.MessagesMessagesClass {
		messages := make([]tg.MessageClass, 0, amount)
		tNow := time.Now()
		for i := 0; i < amount; i++ {
			messages = append(messages, &tg.Message{
				ID:      1,
				Date:    int(tNow.Add(time.Second * time.Duration(i)).Unix()),
				Message: "test" + strconv.Itoa(i),
			})
		}
		messagesClass := tg.MessagesMessagesClass(&tg.MessagesChannelMessages{
			Messages: messages,
		})

		return messagesClass
	}

	answerGpt := func(amount int) []string {
		answer := make([]string, 0, amount)
		for i := 0; i < amount; i++ {
			answer = append(answer, "0")
		}
		return answer
	}

	tests := []struct {
		name      string
		request   *message.ParserRequest
		mocksArgs mocksArgs
		result    result
	}{
		// TODO: Add test cases.
		{
			name: "default: empty bd, new all record, racism, YakiModel",
			request: &message.ParserRequest{
				SourceURL: "https://t.me/maximkatz",
				ScoreType: models.Racism,
				Parser: &message.Parser{
					Type:  models.YakiModel_GPT3_5TURBO,
					Token: "",
				},
				ClientID: "random_string",
			},
			mocksArgs: mocksArgs{
				argumentsTg: []argumentsTg{
					{
						methodName: "ContactSearch",
						args:       []interface{}{"https://t.me/maximkatz"},
						argsReturn: []interface{}{
							&tg.ContactsFound{Chats: []tg.ChatClass{
								&tg.Channel{Username: "maximkat", Title: "Канал Максима"},
								&tg.Channel{Username: "maximkatz", Title: "Канал Максима Каца"},
							}}, nil,
						},
						isCall: true,
					},
					{
						methodName: "MessagesGetHistory",
						args:       []interface{}{&tg.Channel{Username: "maximkatz", Title: "Канал Максима Каца"}, 20, 0},
						argsReturn: []interface{}{messages(20), nil},
						isCall:     true,
					},
				},
				argumentsGpt: []argumentsGpt{
					{
						methodName: "QuestionForGPT",
						args:       []interface{}{},
						argsReturn: []interface{}{answerGpt(20)},
						isCall:     true,
					},
				},
			},
			result: result{
				name:       "Канал Максима Кац",
				typeSource: models.Telegram,
				code:       http.StatusOK,
			},
		},
	}

	srv, srvComponents := mocks.NewMockServer(t)
	rr := httptest.NewRecorder()
	assertTest := assert.New(t)
	index := 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//  setting mock
			for _, arg := range tt.mocksArgs.argumentsTg {
				srvComponents.MockComponents.TgClient.On(arg.methodName, arg.args...).
					Return(arg.argsReturn...).Once()
			}

			argGpt := make([]*openai.ChatCompletionResponse, 0, len(tt.mocksArgs.argumentsGpt))

			for _, arg := range tt.mocksArgs.argumentsGpt {
				for _, argGptChoice := range arg.argsReturn {
					choice, ok := argGptChoice.([]string)
					if !ok {
						t.Fatal("unexpected type of choice")
					}
					for _, argGptChoiceStr := range choice {
						argGpt = append(argGpt, &openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{
								{Message: openai.ChatCompletionMessage{Content: argGptChoiceStr}},
							},
						})

					}

				}
			}

			messagesClass, ok := tt.mocksArgs.argumentsTg[1].argsReturn[0].(tg.MessagesMessagesClass)
			if !ok {
				t.Fatal("unexpected type of message")
			}
			channelMessage, ok := messagesClass.(*tg.MessagesChannelMessages)
			if !ok {
				t.Fatal("unexpected type of message")
			}
			for _, mes := range channelMessage.Messages {
				v, ok := mes.(*tg.Message)
				if !ok {
					t.Fatal("unexpected type of message")
				}
				tt.mocksArgs.argumentsGpt[index].args = append(tt.mocksArgs.argumentsGpt[index].args,
					language_model.AnswerGPTRacism+" "+v.Message)

			}

			for ind, arg := range tt.mocksArgs.argumentsGpt {
				ss, ok := arg.argsReturn[0].([]string)
				if !ok {
					t.Fatal("unexpected type of choice")
				}
				tt.mocksArgs.argumentsGpt[ind].argsReturn = make([]interface{}, 0)
				for _, argGptChoiceStr := range ss {
					tt.mocksArgs.argumentsGpt[ind].argsReturn = append(tt.mocksArgs.argumentsGpt[ind].argsReturn, argGptChoiceStr)

				}
			}

			for _, arg := range tt.mocksArgs.argumentsGpt {
				for i, args := range arg.args {

					srvComponents.MockComponents.GPTClient.On(arg.methodName, args).
						Return(argGpt[i], nil).Once()

				}

			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tt.request)
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/source/parse", &buf)
			if err != nil {
				t.Fatal(err)
			}

			// Использование httptest для записи ответа
			srv.ServeHTTP(rr, req)

			assertTest.Equal(tt.result.code, rr.Code)

			//  check result bd

			ch, ok := tt.mocksArgs.argumentsTg[1].args[0].(*tg.Channel)
			if !ok {
				t.Fatal("unexpected type of channel")
				return
			}
			var res models.SourceDTO
			srvComponents.Components.Db.Select(&models.SourceDTO{}).Find(&res)
			res = models.SourceDTO{
				ID:           1,
				Name:         ch.Username + " " + "@" + ch.Title,
				SourceType:   models.Telegram,
				SourceUrl:    "https://t.me/maximkatz",
				WafflerScore: 0,
				RacismScore:  0,
			}

			index++
		})
	}
}
