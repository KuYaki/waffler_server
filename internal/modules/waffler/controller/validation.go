package controller

import (
	"errors"

	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
)

func ValidatePriceRequest(priceRequest message.PriceRequest) error {
	if priceRequest.SourceUrl == "" {
		return errors.New("Source URL cannot be empty")
	}
	if priceRequest.ScoreType != models.Waffler && priceRequest.ScoreType != models.Racism {
		return errors.New("Invalid score type")
	}
	if priceRequest.Parser.Type != "GPT" && priceRequest.Parser.Type != "YakiModel" {
		return errors.New("Invalid parser type")
	}

	return nil
}
