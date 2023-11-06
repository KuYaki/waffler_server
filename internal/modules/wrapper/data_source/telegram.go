package data_source

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/go-faster/errors"
	tg2 "github.com/gotd/td/tg"
	"log"
	"sort"
	"strings"
	"time"
)

const maxParseOnse = 100

type DataSource struct {
	client telegram.ClientSource
}

type DataSourcer interface {
	ParseChatTelegram(query string, limit int, idMessages []models.RecordDTO) (*DataTelegram, error)
	ContactSearch(query string) (*tg2.Channel, error)
}

func NewDataTelegram(client telegram.ClientSource) DataSourcer {
	return &DataSource{
		client: client,
	}
}

type DataTelegram struct {
	Source  *models.SourceDTO
	Records []models.RecordDTO
}

func (w *DataSource) parseChatTelegram(channel *tg2.Channel, limit int, sessionTs time.Time) ([]models.RecordDTO, error) {
	records := make([]models.RecordDTO, 0, limit)
	for i := 0; i < limit; {
		var limitParser int
		if limit-i > maxParseOnse {
			limitParser = maxParseOnse

		} else {
			limitParser = limit - i

		}

		recordsPart, err := w.parseChat(channel, limitParser, i, sessionTs)
		if err != nil {
			return nil, err
		}
		records = append(records, recordsPart...)

		i += limitParser

	}
	return records, nil
}

func (w *DataSource) ParseChatTelegram(query string, limit int, recordDTOS []models.RecordDTO) (*DataTelegram, error) {
	var newRecords []models.RecordDTO
	var records = make([]models.RecordDTO, 0, limit)
	var sessionTs = time.Now()
	var err error
	channel, err := w.ContactSearch(query)
	if err != nil {
		return nil, err
	}

	if len(recordDTOS) == 0 {
		newRecords, err = w.parseChatTelegram(channel, limit, sessionTs)
		if err != nil {
			return nil, err
		}
		records = append(records, newRecords...)

	} else {
		var messagesTelegramRaw tg2.MessagesMessagesClass
		messagesTelegramRaw, err = w.client.MessagesGetHistoryTime(channel, 100, int(time.Now().Unix()))
		if err != nil {
			return nil, err
		}
		limit -= 100
		if limit > 0 {
			sort.Slice(recordDTOS, func(i, j int) bool {
				return recordDTOS[i].ID > recordDTOS[j].ID
			})

			newRecords, err = messageToRecordDTO(messagesTelegramRaw, sessionTs, channel)
			if err != nil {
				return nil, err
			}

			for recordDTOS[0].CreatedTs.After(newRecords[len(newRecords)-1].CreatedTs) {
				messagesTelegramRaw, err = w.client.MessagesGetHistoryTime(channel, 100, int(newRecords[len(newRecords)-1].CreatedTs.Unix()))
				if err != nil {
					return nil, err
				}
				newRecords, err = messageToRecordDTO(messagesTelegramRaw, sessionTs, channel)
				if err != nil {
					return nil, err
				}
			}
		}

		records = append(records, newRecords...)
	}

	return &DataTelegram{
		Source: &models.SourceDTO{
			Name:       channel.Username + " " + "@" + channel.Title,
			SourceType: models.Telegram,
			SourceUrl:  query,
		},
		Records: records,
	}, nil

}

func (w *DataSource) getMessageForId(channel *tg2.Channel, limit int, AddOffset int) ([]int, error) {
	records := make([]int, 0, limit)
	mes, err := w.client.MessagesGetHistory(channel, limit, AddOffset)
	if err != nil {
		log.Fatalln("failed to get chat:", err)

	}
	res, ok := mes.(*tg2.MessagesChannelMessages)
	if !ok {
		return nil, fmt.Errorf("unknown message type: %T", mes)
	}

	for _, mesRaw := range res.Messages {
		v, ok := mesRaw.(*tg2.Message)
		if ok {
			records = append(records, v.ID)
		} else {
			return nil, fmt.Errorf("unknown message type: %T", v) // ToDo: log
		}

	}

	return records, nil
}

func (w *DataSource) getMessagesForID(channel *tg2.Channel, iDs []int, sessionTs time.Time) ([]models.RecordDTO, error) {
	mes, err := w.client.GetMessagesForID(channel, iDs)
	if err != nil {
		log.Fatalln("failed to get chat:", err)
	}
	return messageToRecordDTO(mes, sessionTs, channel)
}

func (w *DataSource) parseChat(channel *tg2.Channel, limit int, AddOffset int, sessionTs time.Time) ([]models.RecordDTO, error) {
	mes, err := w.client.MessagesGetHistory(channel, limit, AddOffset)
	if err != nil {
		log.Fatalln("failed to get chat:", err)
	}

	return messageToRecordDTO(mes, sessionTs, channel)
}

func messageToRecordDTO(mes tg2.MessagesMessagesClass, sessionTs time.Time, channel *tg2.Channel) ([]models.RecordDTO, error) {
	res, ok := mes.(*tg2.MessagesChannelMessages)
	if !ok {
		return nil, fmt.Errorf("unknown message type: %T", mes)
	}
	records := make([]models.RecordDTO, 0, len(res.Messages))

	for _, mesRaw := range res.Messages {
		v, ok := mesRaw.(*tg2.Message)
		if ok {
			records = append(records,
				models.RecordDTO{
					RecordText: v.Message,
					CreatedTs:  time.Unix(int64(v.Date), 0),
					SessionTs:  sessionTs,
					RecordURL:  fmt.Sprintf("https://t.me/%s/%d", channel.Username, v.ID),
				})
		} else {
			return nil, fmt.Errorf("unknown message type: %T", v)
		}

	}
	return records, nil
}

func (w *DataSource) ContactSearch(query string) (*tg2.Channel, error) {
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
