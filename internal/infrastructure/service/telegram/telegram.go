package telegram

import (
	"context"
	"fmt"
	"github.com/KuYaki/waffler_server/config"
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

func (t *Telegram) ParseChat(query string, limit int) {
	f, err := t.Client.ContactsSearch(context.Background(), &tg2.ContactsSearchRequest{
		Q:     query,
		Limit: limit,
	})
	if err != nil {
		log.Fatalln("failed to get contacts:", err)
	}

	targetUsername := strings.Split(query, "/")

	username := targetUsername[len(targetUsername)-1]
	var ss *tg2.Channel
	for _, chat := range f.GetChats() {
		ss = chat.(*tg2.Channel)
		if ss.Username == username {
			break
		}
	}
	fmt.Println(ss.Username)

	mes, err := t.Client.MessagesGetHistory(context.Background(), &tg2.MessagesGetHistoryRequest{
		Peer: &tg2.InputPeerChannel{
			ChannelID:  ss.ID,
			AccessHash: ss.AccessHash},
		Limit: limit,
	})
	if err != nil {
		log.Fatalln("failed to get chat:", err)
	}
	res := mes.(*tg2.MessagesChannelMessages)

	log.Println(res.Messages)
	//return res.Messages
}
