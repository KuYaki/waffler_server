package storage

import (
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/go-faster/errors"
	"gorm.io/gorm"
)

type BotStorager interface {
	CreateWebhook(webhook *models.WebhookDTO) error
	UpdateWebhook(webhook *models.WebhookDTO) error
	TakeWebhook(webhookName string) (*models.WebhookDTO, error)
	ExistWebhook(webhookName string) (bool, error)
	DeleteWebhook(webhook *models.WebhookDTO) error
}

type BotStorage struct {
	conn *gorm.DB
}

func NewWafflerStorage(conn *gorm.DB) BotStorager {
	return &BotStorage{conn: conn}
}

func (b *BotStorage) CreateWebhook(webhook *models.WebhookDTO) error {
	return b.conn.Create(webhook).Error
}

func (b *BotStorage) TakeWebhook(webhookName string) (*models.WebhookDTO, error) {
	var webhookDTO models.WebhookDTO
	err := b.conn.Where("name = ?", webhookName).First(&webhookDTO).Error
	if err != nil {
		return nil, err
	}
	return &webhookDTO, nil

}

func (b *BotStorage) UpdateWebhook(webhook *models.WebhookDTO) error {
	return b.conn.Where("name = ?", webhook.Name).Updates(webhook).Error
}

func (b *BotStorage) ExistWebhook(webhookName string) (bool, error) {

	res := b.conn.Where("name = ?", webhookName).First(&models.WebhookDTO{})
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) && res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil

}

func (b *BotStorage) DeleteWebhook(webhook *models.WebhookDTO) error {
	return b.conn.Delete(webhook).Error
}
