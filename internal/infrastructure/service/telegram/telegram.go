package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/models"
	tg2 "github.com/gotd/td/tg"
	tgpro "github.com/jaskaur18/gotgproto"
	"github.com/jaskaur18/gotgproto/sessionMaker"
	"log"
	"strings"
	"time"
)

type Telegram struct {
	Client *tg2.Client
}

const (
	maxParseOnse       = 100
	limitContactSearch = 100
)

func NewTelegram(conf *config.Telegram) (*Telegram, error) {
	clientType := tgpro.ClientType{
		Phone: conf.Phone,
	}

	client, err := tgpro.NewClient(
		// Get AppID from https://my.telegram.org/apps
		conf.AppID,
		// Get ApiHash from https://my.telegram.org/apps
		conf.ApiHash,
		// ClientType, as we defined above
		clientType,
		// Optional parameters of client
		&tgpro.ClientOpts{
			Session: sessionMaker.NewSession(conf.Phone, sessionMaker.Session, sessionMaker.NewSessionOpts{
				SessionName: conf.Phone,
				SessionPath: "./telegram_sessions",
			}),
		},
	)
	if err != nil {
		return nil, err
	}
	tg := client.API()
	return &Telegram{
		Client: tg,
	}, nil
}

type UserData struct {
	Id           int
	Name         string
	SoursType    string
	SoursUrl     string
	WafflerScore int
	Messages     []MessageData
}

type MessageData struct {
	Id          int
	RecordText  string
	Score       int
	TimeMessage time.Time
}

func (t *Telegram) ContactSearch(query string) (*tg2.Channel, error) {
	f, err := t.Client.ContactsSearch(context.Background(), &tg2.ContactsSearchRequest{
		Q:     query,
		Limit: limitContactSearch,
	})
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

func (t *Telegram) ParseChatTelegram(query string, limit int) (*DataTelegram, error) {
	channel, err := t.ContactSearch(query)
	if err != nil {
		return nil, err
	}

	dataTg := &DataTelegram{
		&models.SourceDTO{
			Name:       channel.Username + " " + "@" + channel.Title,
			SourceType: models.Telegram,
			SourceUrl:  query,
		},
		make([]*models.RecordDTO, 0, limit),
	}

	for i := 0; i < limit; {
		var limitParser int
		if limit-i > maxParseOnse {
			limitParser = maxParseOnse

		} else {
			limitParser = limit - i

		}

		err = t.ParseChat(dataTg, channel, limitParser, i)
		if err != nil {
			return nil, err
		}

		i += limitParser

	}
	return dataTg, nil

}
func (t *Telegram) ParseChat(dataTg *DataTelegram, channel *tg2.Channel, limit int, AddOffset int) error {

	mes, err := t.Client.MessagesGetHistory(context.Background(), &tg2.MessagesGetHistoryRequest{
		Peer: &tg2.InputPeerChannel{
			ChannelID:  channel.ID,
			AccessHash: channel.AccessHash},
		Limit:     limit,
		AddOffset: AddOffset,
	})
	if err != nil {
		log.Fatalln("failed to get chat:", err)
	}
	res := mes.(*tg2.MessagesChannelMessages) //  ToDo: switch type

	for _, mesRaw := range res.Messages {
		switch v := mesRaw.(type) {
		case *tg2.MessageEmpty: // messageEmpty#90a6ca84
		case *tg2.Message: // message#38116ee0
			message := mesRaw.(*tg2.Message)
			dataTg.Records = append(dataTg.Records,
				&models.RecordDTO{
					RecordText: message.Message,
					CreatedAt:  time.Unix(int64(message.Date), 0),
				})
		case *tg2.MessageService: // messageService#2b085862
		default:
			return fmt.Errorf("unknown message type: %T", v) // ToDo: log
		}

	}

	return nil
}

type DataTelegram struct {
	Source  *models.SourceDTO
	Records []*models.RecordDTO
}
