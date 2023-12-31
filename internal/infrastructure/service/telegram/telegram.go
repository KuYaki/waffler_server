package telegram

import (
	"context"
	"github.com/KuYaki/waffler_server/config"
	tg2 "github.com/gotd/td/tg"
	tgpro "github.com/jaskaur18/gotgproto"
	"github.com/jaskaur18/gotgproto/sessionMaker"
	"time"
)

type Telegram struct {
	Client *tg2.Client
}

const (
	limitContactSearch = 100
)

type ClientSource interface {
	ContactSearch(query string) (*tg2.ContactsFound, error)
	MessagesGetHistory(channel *tg2.Channel, limit int, AddOffset int) (tg2.MessagesMessagesClass, error)
}

func NewTelegram(conf *config.Telegram) (ClientSource, error) {
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

func (t *Telegram) ContactSearch(query string) (*tg2.ContactsFound, error) {
	return t.Client.ContactsSearch(context.Background(), &tg2.ContactsSearchRequest{
		Q:     query,
		Limit: limitContactSearch,
	})
}

func (t *Telegram) MessagesGetHistory(channel *tg2.Channel, limit int, AddOffset int) (tg2.MessagesMessagesClass, error) {
	return t.Client.MessagesGetHistory(context.Background(), &tg2.MessagesGetHistoryRequest{
		Peer: &tg2.InputPeerChannel{
			ChannelID:  channel.ID,
			AccessHash: channel.AccessHash},
		Limit:     limit,
		AddOffset: AddOffset,
	})

}
