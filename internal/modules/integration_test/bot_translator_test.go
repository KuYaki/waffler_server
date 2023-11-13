//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator"
	"github.com/KuYaki/waffler_server/mocks"
	"github.com/goccy/go-json"
	"github.com/gotd/td/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetWebhook(t *testing.T) {
	type sender struct {
		update      bot_translator.Update
		host        string
		returnError error
	}

	type contactSearch struct {
		arguments string
		channel   []tg.ChatClass
		error     error
		isCall    bool
	}
	type messagesGetHistoryTime struct {
		ch     *tg.Channel
		limit  int
		offset int
		time   interface{}
		error  error
		isCall bool

		returnMock tg.MessagesMessagesClass
	}

	type mocksArgs struct {
		contactSearch          contactSearch
		messagesGetHistoryTime []messagesGetHistoryTime
		sender                 []sender
	}

	type result struct {
		name       string
		typeSource models.SourceType
		code       int
	}

	channels := []*tg.Channel{
		{
			Username:   "maximkatz",
			Title:      "Канал Максима Каца",
			ID:         1,
			AccessHash: 1,
		},
	}

	tests := []struct {
		name          string
		request       bot_translator.SetWebhookParams
		requestDelete bot_translator.DeleteWebhookParams
		mocksArgs     mocksArgs
		result        result
	}{
		{
			name:          "default telegram",
			request:       bot_translator.SetWebhookParams{URL: "https://t.me/maximkatz"},
			requestDelete: bot_translator.DeleteWebhookParams{URL: "https://t.me/maximkatz"},
			mocksArgs: mocksArgs{
				contactSearch: contactSearch{
					arguments: "https://t.me/maximkatz",
					channel:   []tg.ChatClass{channels[0]},
					error:     nil,
					isCall:    true,
				},
				messagesGetHistoryTime: []messagesGetHistoryTime{{
					ch:     channels[0],
					limit:  1,
					offset: 0,
					time:   mock.Anything,
					error:  nil,
					isCall: true,
					returnMock: &tg.MessagesChannelMessages{
						Messages: []tg.MessageClass{
							&tg.Message{
								ID: 1000,
							},
						},
					},
				},

					{
						ch:     channels[0],
						limit:  10,
						offset: 0,
						time:   mock.Anything,
						error:  nil,
						isCall: true,
						returnMock: &tg.MessagesChannelMessages{
							Messages: []tg.MessageClass{
								&tg.Message{
									ID: 1000,
								},
								&tg.Message{
									ID: 1001,
								},
								&tg.Message{
									ID: 1002,
								},
								&tg.Message{
									ID: 1003,
								},
								&tg.Message{
									ID: 1004,
								},
								&tg.Message{
									ID: 1005,
								},
								&tg.Message{
									ID: 1006,
								},
								&tg.Message{
									ID: 1007,
								},
								&tg.Message{
									ID: 1008,
								},
								&tg.Message{
									ID: 1009,
								},
							},
						},
					},

					{
						ch:     channels[0],
						limit:  10,
						offset: 0,
						time:   mock.Anything,
						error:  nil,
						isCall: true,
						returnMock: &tg.MessagesChannelMessages{
							Messages: []tg.MessageClass{
								&tg.Message{
									ID: 1003,
								},
								&tg.Message{
									ID: 1004,
								},
								&tg.Message{
									ID: 1005,
								},
								&tg.Message{
									ID: 1006,
								},
								&tg.Message{
									ID: 1007,
								},
								&tg.Message{
									ID: 1008,
								},
								&tg.Message{
									ID: 1009,
								},
								&tg.Message{
									ID: 1010,
								},
								&tg.Message{
									ID: 1011,
								},
								&tg.Message{
									ID: 1012,
								},
							},
						},
					},
					{
						ch:     channels[0],
						limit:  10,
						offset: 0,
						time:   mock.Anything,
						error:  nil,
						isCall: true,
						returnMock: &tg.MessagesChannelMessages{
							Messages: []tg.MessageClass{
								&tg.Message{
									ID: 1003,
								},
								&tg.Message{
									ID: 1004,
								},
								&tg.Message{
									ID: 1005,
								},
								&tg.Message{
									ID: 1006,
								},
								&tg.Message{
									ID: 1007,
								},
								&tg.Message{
									ID: 1008,
								},
								&tg.Message{
									ID: 1009,
								},
								&tg.Message{
									ID: 1010,
								},
								&tg.Message{
									ID: 1011,
								},
								&tg.Message{
									ID: 1012,
								},
							},
						},
					},
				},
				sender: []sender{{
					update: bot_translator.Update{
						ID: 1,
						ChannelPost: &tg.Message{
							ID: 1001,
						},
					},
				},

					{
						update: bot_translator.Update{
							ID: 2,
							ChannelPost: &tg.Message{
								ID: 1002,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 3,
							ChannelPost: &tg.Message{
								ID: 1003,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 4,
							ChannelPost: &tg.Message{
								ID: 1004,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 5,
							ChannelPost: &tg.Message{
								ID: 1005,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 6,
							ChannelPost: &tg.Message{
								ID: 1006,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 7,
							ChannelPost: &tg.Message{
								ID: 1007,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 8,
							ChannelPost: &tg.Message{
								ID: 1008,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 9,
							ChannelPost: &tg.Message{
								ID: 1009,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 10,
							ChannelPost: &tg.Message{
								ID: 1010,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 11,
							ChannelPost: &tg.Message{
								ID: 1011,
							},
						},
					},
					{
						update: bot_translator.Update{
							ID: 12,
							ChannelPost: &tg.Message{
								ID: 1012,
							},
						},
					},
				},
			},

			result: result{name: "Канал Максима Каца",
				typeSource: models.Telegram,
				code:       http.StatusOK},
		},
	}

	tests[0].mocksArgs.messagesGetHistoryTime[0].time = time.Now()

	srv, srvComponents := mocks.NewMockServer(t)
	rr := httptest.NewRecorder()
	assertTest := assert.New(t)
	i := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.mocksArgs.contactSearch.isCall {
				srvComponents.MockComponents.TgClient.On("ContactSearch", tt.mocksArgs.contactSearch.arguments).
					Return(&tg.ContactsFound{Chats: tt.mocksArgs.contactSearch.channel}, nil).Once()
			}

			if tt.mocksArgs.messagesGetHistoryTime[i].isCall {
				for _, val := range tt.mocksArgs.messagesGetHistoryTime {
					tempMockArgs := val
					srvComponents.MockComponents.TgClient.On("MessagesGetHistoryTime", tempMockArgs.ch,
						tempMockArgs.limit, tempMockArgs.offset, mock.Anything).
						Return(tempMockArgs.returnMock, tempMockArgs.error).Once()
				}
			}

			for _, val := range tt.mocksArgs.sender {
				tempMockArgs := val
				tempMockArgs.host = "test.com"
				srvComponents.MockComponents.SenderWebhook.On("SendUpdate", tempMockArgs.update, tempMockArgs.host).
					Return(val.returnError).Once()

			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tt.request)
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/bot_translator/set_webhook", &buf)
			if err != nil {
				t.Fatal(err)
			}
			req.Host = "test.com"

			// Использование httptest для записи ответа
			srv.ServeHTTP(rr, req)

			assertTest.Equal(tt.result.code, rr.Code)
			time.Sleep(time.Second * 35)

			err = json.NewEncoder(&buf).Encode(tt.requestDelete)
			if err != nil {
				log.Fatal(err)
			}

			req, err = http.NewRequest("POST", "/bot_translator/delete_webhook", &buf)
			if err != nil {
				t.Fatal(err)
			}
			req.Host = "test.com"
			srv.ServeHTTP(rr, req)

			assertTest.Equal(tt.result.code, rr.Code)
			var webhookDTO models.WebhookDTO
			res := srvComponents.Components.Db.First(&webhookDTO).
				Where("webhook_url = ?", tt.request.URL)
			if res.Error != nil {
				t.Fatal(res.Error)
			}
			assertTest.Equal(*webhookDTO.Active, false)

			i++
		})
	}
}
