package service

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/webhook"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator"
	"github.com/KuYaki/waffler_server/internal/modules/bot_translator/storage"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"time"
)

type BotTranslatorServicer interface {
	SetWebhook(setWebhookParams bot_translator.SetWebhookParams, host string) error
	DeleteWebhook(deleteWebhookParams bot_translator.DeleteWebhookParams, host string) error
}

type Service struct {
	tgClient   telegram.ClientSource
	tgWrapper  data_source.DataSourcer
	sender     webhook.SenderWebhooker
	botStorage storage.BotStorager
	log        *zap.Logger

	workers map[string]*worker
}

func NewWafflerService(storage storage.BotStorager, components *component.Components) BotTranslatorServicer {
	return &Service{
		tgClient:   components.TgClient,
		tgWrapper:  components.TgWrapper,
		sender:     components.SenderWebhook,
		botStorage: storage,
		log:        components.Logger,
		workers:    make(map[string]*worker),
	}
}

func (s *Service) DeleteWebhook(deleteWebhookParams bot_translator.DeleteWebhookParams, host string) error {
	s.workers[host+deleteWebhookParams.URL].stop()
	return nil
}

func (s *Service) SetWebhook(setWebhookParams bot_translator.SetWebhookParams, host string) error {
	ch, err := s.tgWrapper.ContactSearch(setWebhookParams.URL)
	if err != nil {
		return err
	}

	var timeInit = time.Now()

	exist, err := s.botStorage.ExistWebhook(host + setWebhookParams.URL)
	if err != nil {
		return err
	}
	mes, err := s.tgClient.MessagesGetHistoryTime(ch, 1, 0, int(time.Now().Unix()))
	if err != nil {
		s.log.Error("error: get messages", zap.Error(err))
		return err
	}

	res, ok := mes.(*tg.MessagesChannelMessages)
	if !ok {
		return fmt.Errorf("unknown message type: %T", mes)
	}

	v, ok := res.Messages[0].(*tg.Message)
	if !ok {
		return fmt.Errorf("unknown message type: %T", res.Messages[0])
	}

	if !exist {
		err := s.botStorage.CreateWebhook(&models.WebhookDTO{
			Name:            host + setWebhookParams.URL,
			LastTimeRequest: timeInit,
			LastIdMessage:   v.ID,
			UpdateId:        1,
		})
		if err != nil {
			return err
		}
	}

	webhookDTO, err := s.botStorage.TakeWebhook(host + setWebhookParams.URL)
	if err != nil {
		return err
	}

	if setWebhookParams.DropPendingUpdates {
		webhookDTO.LastTimeRequest = timeInit
		webhookDTO.LastIdMessage = v.ID
	}

	worker := s.createWorker(setWebhookParams.URL, host, ch, webhookDTO.LastTimeRequest,
		webhookDTO.LastIdMessage, webhookDTO.UpdateId)

	s.workers[worker.host+setWebhookParams.URL] = worker
	go s.workers[worker.host+setWebhookParams.URL].run()

	return nil
}

type worker struct {
	updateId      int
	lastIdMessage int
	host          string
	URLResponse   string
	sender        webhook.SenderWebhooker
	botStorage    storage.BotStorager
	scanChan      *tg.Channel
	tgClient      telegram.ClientSource
	time          time.Time
	signalStop    chan struct{}
	log           *zap.Logger
}

func (s *Service) createWorker(URLResponse, host string, channel *tg.Channel, timeInit time.Time, lastIdMessage, updateId int) *worker {
	return &worker{
		updateId:      updateId,
		lastIdMessage: lastIdMessage,
		host:          host,
		URLResponse:   URLResponse,
		sender:        s.sender,
		botStorage:    s.botStorage,
		scanChan:      channel,
		tgClient:      s.tgClient,
		time:          timeInit,
		signalStop:    make(chan struct{}),
		log:           s.log,
	}

}

const timeRequestInterval = 10 * time.Second

func (w *worker) run() {
	var err error
	activeWebhook := true
	err = w.botStorage.UpdateWebhook(&models.WebhookDTO{
		Name:   w.host + w.URLResponse,
		Active: &activeWebhook,
	})
	if err != nil {
		w.log.Error("error: update webhook", zap.Error(err))
	}
	ticker := time.NewTicker(timeRequestInterval)
	stopSignal := false

	for {
		select {
		case <-ticker.C:
			res := w.getNewMessages()
			for _, mesRaw := range res {
				//time.Sleep(time.Second * 1)

				upd := bot_translator.Update{
					ID:          w.updateId,
					ChannelPost: mesRaw,
				}
				err := w.sender.SendUpdate(upd, w.host)
				if err != nil {
					break
				}
				w.updateId++
			}
			w.time.Add(timeRequestInterval)
			if err != nil {
				break
			}
			if res != nil {
				w.lastIdMessage = res[len(res)-1].ID
			}

		case <-w.signalStop:
			stopSignal = true
			break
		}
		if err != nil || stopSignal {
			break
		}
	}
	if err != nil {
		w.log.Error("error: worker", zap.Error(err))
	}
	activeWebhook = false
	err = w.botStorage.UpdateWebhook(&models.WebhookDTO{
		Name:            w.host + w.URLResponse,
		LastTimeRequest: w.time,
		LastIdMessage:   w.lastIdMessage,
		UpdateId:        w.updateId,
		Active:          &activeWebhook,
	})
	if err != nil {
		w.log.Error("error: db", zap.Error(err))
	}
}

func (w *worker) getNewMessages() []*tg.Message {
	fmt.Println(w.time.Unix())
	mes, err := w.tgClient.MessagesGetHistoryTime(w.scanChan, 10, 0, int(w.time.Unix()))
	if err != nil {
		w.log.Error("error: get messages", zap.Error(err))
	}

	mess, err := w.getNewMessagesHelper(mes)
	lastIDMess := w.lastIdMessage
	if mess != nil {
		lastIDMess = mess[len(mess)-1].ID
	}

	for lastIDMess < w.lastIdMessage {
		mes, err := w.tgClient.MessagesGetHistoryTime(w.scanChan, 100, 0, mess[len(mess)-1].GetDate())
		if err != nil {
			w.log.Error("error: get messages", zap.Error(err))
		}

		messPart, err := w.getNewMessagesHelper(mes)
		if err != nil {
			w.log.Error("error: get messages", zap.Error(err))
		}

		mess = append(mess, messPart...)
	}

	return mess
}

func (w *worker) getNewMessagesHelper(mes tg.MessagesMessagesClass) ([]*tg.Message, error) {
	res, ok := mes.(*tg.MessagesChannelMessages)
	if !ok {
		return nil, fmt.Errorf("unknown message type: %T", mes)
	}
	var result []*tg.Message
	for _, mesRaw := range res.Messages {
		//time.Sleep(time.Second * 1)
		v, ok := mesRaw.(*tg.Message)
		if !ok {
			return nil, fmt.Errorf("unknown message type: %T", mesRaw)
		}
		if v.ID > w.lastIdMessage {
			result = append(result, v)
		}

	}
	return result, nil

}

func (w *worker) stop() {
	w.signalStop <- struct{}{}
}
