package component

import (
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/responder"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Conf         config.AppConf
	TokenManager cryptography.TokenManager
	Responder    responder.Responder
	Decoder      godecoder.Decoder
	Logger       *zap.Logger
	Hash         cryptography.Hasher
	Tg           *telegram.Telegram
	Gpt          *gpt.ChatGPT
}

func NewComponents(conf config.AppConf, tokenManager cryptography.TokenManager, responder responder.Responder, decoder godecoder.Decoder, hash cryptography.Hasher, telegram *telegram.Telegram, logger *zap.Logger, gpt *gpt.ChatGPT) *Components {
	return &Components{Conf: conf,
		TokenManager: tokenManager,
		Responder:    responder,
		Decoder:      decoder,
		Hash:         hash,
		Logger:       logger,
		Tg:           telegram,
		Gpt:          gpt}
}
