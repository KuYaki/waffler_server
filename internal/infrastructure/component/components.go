package component

import (
	"github.com/KuYaki/waffler_server/config"
	middleware "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/webhook"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/language_model"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Conf          *config.AppConf
	TokenManager  cryptography.TokenManager
	Token         *middleware.Token
	Responder     responder.Responder
	Decoder       godecoder.Decoder
	Logger        *zap.Logger
	Hash          cryptography.Hasher
	TgWrapper     data_source.DataSourcer
	TgClient      telegram.ClientSource
	Gpt           language_model.LanguageModel
	SenderWebhook webhook.SenderWebhooker
}

func NewComponents(conf *config.AppConf, tokenManager cryptography.TokenManager, token *middleware.Token,
	responder responder.Responder, decoder godecoder.Decoder, hash cryptography.Hasher, dataSource data_source.DataSourcer,
	logger *zap.Logger, gpt language_model.LanguageModel, tgClient telegram.ClientSource, SenderWebhook webhook.SenderWebhooker) *Components {
	return &Components{Conf: conf,
		TokenManager:  tokenManager,
		Token:         token,
		Responder:     responder,
		Decoder:       decoder,
		Hash:          hash,
		Logger:        logger,
		TgWrapper:     dataSource,
		TgClient:      tgClient,
		Gpt:           gpt,
		SenderWebhook: SenderWebhook}
}
