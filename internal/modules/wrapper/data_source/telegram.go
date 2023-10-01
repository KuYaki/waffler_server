package data_source

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/go-faster/errors"
	tg2 "github.com/gotd/td/tg"
	"log"
	"strings"
	"time"
)

const maxParseOnse = 100

type DataSource struct {
	Source  *models.SourceDTO
	Records []*models.RecordDTO
	client  telegram.ClientSource
}

type DataSourcer interface {
	ParseChatTelegramOld(channel *tg2.Channel, limit int) error
	ContactSearchOld(query string) (*tg2.Channel, error)
}

func NewDataTelegram(sourceUrl string) *DataSource {
	return &DataSource{
		Source: &models.SourceDTO{
			SourceType: models.Telegram,
			SourceUrl:  sourceUrl,
		},
		Records: make([]*models.RecordDTO, 0),
	}
}

func (w *DataSource) ParseChatTelegramOld(channel *tg2.Channel, limit int) error {
	records := make([]*models.RecordDTO, 0, limit)

	for i := 0; i < limit; {
		var limitParser int
		if limit-i > maxParseOnse {
			limitParser = maxParseOnse

		} else {
			limitParser = limit - i

		}

		recordsPart, err := w.parseChatOld(channel, limitParser, i)
		if err != nil {
			return err
		}
		records = append(records, recordsPart...)

		i += limitParser

	}
	w.Records = records
	return nil

}

func (w *DataSource) ContactSearchOld(query string) (*tg2.Channel, error) {
	f, err := w.client.ContactSearch(query)
	if err != nil {
		return nil, err
	}
	urlSplit := strings.Split(query, "/")

	usernameTarget := urlSplit[len(urlSplit)-1]

	var channel *tg2.Channel
	var found bool

	chats := f.GetChats()
	for _, chat := range chats {
		channel = chat.(*tg2.Channel)
		if channel.Username == usernameTarget {
			found = true
		}
		if found {
			break
		}

	}

	if !found {
		for _, chat := range chats {
			channel = chat.(*tg2.Channel)
			if strings.Contains(channel.Username, usernameTarget) {
				found = true
			}
			if found {
				break
			}

		}
	}

	if !found {
		return nil, errors.New("not found")
	}

	return channel, nil
}

func (w *DataSource) parseChatOld(channel *tg2.Channel, limit int, AddOffset int) ([]*models.RecordDTO, error) {
	records := make([]*models.RecordDTO, 0, limit)
	mes, err := w.client.MessagesGetHistory(channel, limit, AddOffset)
	if err != nil {
		log.Fatalln("failed to get chat:", err)
	}
	res := mes.(*tg2.MessagesChannelMessages) //  ToDo: switch type

	for _, mesRaw := range res.Messages {
		switch v := mesRaw.(type) {
		case *tg2.MessageEmpty: // messageEmpty#90a6ca84
		case *tg2.Message: // message#38116ee0
			message := mesRaw.(*tg2.Message)
			records = append(records,
				&models.RecordDTO{
					RecordText: message.Message,
					CreatedAt:  time.Unix(int64(message.Date), 0),
				})
		case *tg2.MessageService: // messageService#2b085862
		default:
			return nil, fmt.Errorf("unknown message type: %T", v) // ToDo: log
		}

	}

	return records, nil
}
