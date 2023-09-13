package component

import (
	"github.com/KuYaki/waffler_server/config"
	middleware "github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Conf         config.AppConf
	TokenManager cryptography.TokenManager
	Token        *middleware.Token
	Responder    responder.Responder
	Decoder      godecoder.Decoder
	Logger       *zap.Logger
	Hash         cryptography.Hasher
	Tg           *telegram.Telegram
}

func NewComponents(conf config.AppConf, tokenManager cryptography.TokenManager, token *middleware.Token, responder responder.Responder, decoder godecoder.Decoder, hash cryptography.Hasher, telegram *telegram.Telegram, logger *zap.Logger) *Components {
	return &Components{Conf: conf,
		TokenManager: tokenManager,
		Token:        token,
		Responder:    responder,
		Decoder:      decoder,
		Hash:         hash,
		Logger:       logger,
		Tg:           telegram}
}
