//go:build integration
// +build integration

package integration_test_test

import (
	"fmt"
	middleware "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	uservice "github.com/KuYaki/waffler_server/internal/modules/user/service"
	wservice "github.com/KuYaki/waffler_server/internal/modules/waffler/service"
	"github.com/KuYaki/waffler_server/mocks"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWaffl_Search(t *testing.T) {
	type fields struct {
		service     wservice.WafflerServicer
		log         *zap.Logger
		token       *middleware.Token
		userService uservice.Userer
		Responder   responder.Responder
		Decoder     godecoder.Decoder
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

	srv := mocks.MockServer(t)

	req, err := http.NewRequest("GET", "/", nil)
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
