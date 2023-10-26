//go:build integration
// +build integration

package integration_test_test

import (
	"bytes"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/mocks"
	"github.com/goccy/go-json"
	"github.com/gotd/td/tg"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mocksArgs.isCall {
				srvComponents.TgClient.On("ContactSearch", tt.mocksArgs.arguments).
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
			assertTest.Equal(tt.result.name, res.Name)
			assertTest.Equal(tt.result.typeSource, res.Type)
			assertTest.Equal(tt.result.code, rr.Code)

		})
	}
}
