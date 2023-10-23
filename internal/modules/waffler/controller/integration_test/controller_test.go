//go:build integration
// +build integration

package integration_test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/mocks"
	"github.com/gotd/td/tg"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWaffl_Search(t *testing.T) {
	type fields struct {
		mocksArgs []mocks.Args
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}

	srv, srvComponents := mocks.NewMockServer(t)
	ret := &tg.ContactsFound{Chats: []tg.ChatClass{&tg.Channel{Username: "maximkatz", Title: "Канал Максима Каца"}}}
	srvComponents.TgClient.On("ContactSearch", "https://t.me/maximkatz").Return(ret, nil).Once()

	body := message.SourceURL{SourceUrl: "https://t.me/maximkatz"}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/source/info", &buf)
	if err != nil {
		t.Fatal(err)
	}

	// Использование httptest для записи ответа
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
