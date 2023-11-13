package bot_translator

import (
	"github.com/go-telegram/bot/models"
	"github.com/gotd/td/tg"
)

type SetWebhookParams struct {
	URL                string           `json:"url"`
	Certificate        models.InputFile `json:"certificate,omitempty"`
	IPAddress          string           `json:"ip_address,omitempty"`
	MaxConnections     int              `json:"max_connections,omitempty"`
	AllowedUpdates     []string         `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool             `json:"drop_pending_updates,omitempty"`
	SecretToken        string           `json:"secret_token,omitempty"`
}

type Update struct {
	ID          int         `json:"update_id"`
	ChannelPost *tg.Message `json:"channel_post,omitempty"`
}

type DeleteWebhookParams struct {
	URL string `json:"url"`
}
