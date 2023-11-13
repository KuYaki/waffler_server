package webhook

import (
	"bytes"
	"fmt"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator"
	"github.com/goccy/go-json"
	"net/http"
)

type webhookSender struct {
}

type SenderWebhooker interface {
	SendUpdate(upd bot_translator.Update, host string) error
}

func NewWebhookSender() SenderWebhooker {
	return &webhookSender{}
}

func (w *webhookSender) SendUpdate(upd bot_translator.Update, host string) error {
	payload, err := json.Marshal(upd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, host, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}
