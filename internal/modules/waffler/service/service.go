package service

import (
	"fmt"
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/gpt"
	"net/url"
	"sync"

	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/language_model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type WafflerServicer interface {
	Search(search *message.Search) (*message.SearchResponse, error)
	InfoSource(urlSearch string) (*message.InfoRequest, error)
	Score(request *message.ScoreRequest) (*message.ScoreResponse, error)
	ParseSource(search *message.ParserRequest, updateChan chan<- bool) error
	PriceSource(request *message.PriceRequest) (*message.PriceResponse, error)
}

type WafflerService struct {
	storage storage.WafflerStorager
	log     *zap.Logger
	tg      data_source.DataSourcer
	gpt     language_model.LanguageModel
}

func NewWafflerService(storage storage.WafflerStorager, components *component.Components) *WafflerService {
	return &WafflerService{storage: storage, log: components.Logger, tg: components.Tg, gpt: components.Gpt}
}

func (w *WafflerService) Score(request *message.ScoreRequest) (*message.ScoreResponse, error) {

	scoreResponse := &message.ScoreResponse{
		Cursor: &message.Cursor{
			Offset: request.Cursor.Offset,
		},
		Records: []message.Record{},
	}

	orders, err := convertRecordOrder(request.Order)
	if err != nil {
		return nil, err
	}

	records, err := w.storage.ListRecordsSourceID(request.SourceId)
	if err != nil {
		return nil, err
	}

	var racismRecords []models.RacismDTO

	switch request.Type {
	case models.Racism:
		racismRecords, err = w.storage.ListRacismRecordsSourceIDCursor(request.SourceId, orders, request.Cursor.Offset, request.Limit)
		if err != nil {
			return nil, err
		}
	}

	if len(racismRecords) == 0 {
		scoreResponse.Cursor = nil
	} else {
		scoreResponse.Cursor.Offset += len(records)

		scoreResponse.Records = make([]message.Record, 0, len(records))
		for i := range racismRecords {
			scoreResponse.Records = append(scoreResponse.Records, message.Record{
				RecordText: records[i].RecordText,
				Score:      int(racismRecords[i].Score.Int64),
				Timestamp:  racismRecords[i].CreatedTs,
			})
		}
	}

	return scoreResponse, nil
}

func (w *WafflerService) InfoSource(urlSearch string) (*message.InfoRequest, error) {
	res := &message.InfoRequest{}

	urlParse, err := url.Parse(urlSearch)
	if err != nil {
		return nil, err
	}

	switch urlParse.Host {
	case "t.me":
		channel, err := w.tg.ContactSearch(urlSearch)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return nil, err
		}

		res.Name = channel.Title
		res.Type = models.Telegram
	default:
		res.Name = "Unknown"
		res.Type = models.Unknown
	}

	return res, nil
}

func (w *WafflerService) parseSourceTypeRacism(search *message.ParserRequest, dataTelegram *data_source.DataTelegram, updateChan chan<- bool) error {
	g := errgroup.Group{}
	g.SetLimit(20)

	listRascismRecords, err := w.storage.ListRacismRecords(&models.RacismDTO{SourceID: dataTelegram.Source.ID})
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	newRacismRecords := make([]models.RacismDTO, len(dataTelegram.Records))
	lanModel := w.createLanModel(search)
	for i, r := range dataTelegram.Records {
		i, r := i, r
		existRasism := false
		for _, rasR := range listRascismRecords {
			if r.ID == rasR.RecordID {
				existRasism = true
				break
			}
		}
		if existRasism {
			updateChan <- true
			continue
		}
		g.Go(func() error {
			var err error
			res, err := lanModel.ConstructQuestionGPT(r.RecordText, search.ScoreType)
			if res != nil {
				score := models.NewNullInt64(int64(*res))
				if score.Int64 < 0 {
					score = models.NullInt64{}
				}
				newRacismRecords[i] = models.RacismDTO{
					Score:      score,
					ParserType: models.GPT3_5TURBO,
					CreatedTs:  r.CreatedTs,
					RecordID:   r.ID,
					SourceID:   dataTelegram.Source.ID,
				}
				updateChan <- true
				return nil
			} else {
				if err != nil {
					w.log.Error("error: search", zap.Error(err))
					return err
				}
			}

			return nil
		})

	}
	err = g.Wait()
	if err != nil {
		w.log.Warn("error: search", zap.Error(err))
		return err
	}

	source, err := w.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	err = w.storage.CreateRacismRecords(newRacismRecords)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	if len(dataTelegram.Records) != 0 {

		racismRecords, err := w.storage.ListRacismRecords(&models.RacismDTO{SourceID: dataTelegram.Source.ID})
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}

		source.RacismScore = averageRacismScore(racismRecords)

		err = w.storage.UpdateSource(source)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	return err
}

func (w *WafflerService) ParseSource(search *message.ParserRequest, updateChan chan<- bool) error {
	var limit = 20 // TODO: search.Limit
	var source *models.SourceDTO
	var records []models.RecordDTO
	var err error
	source, err = w.storage.SearchBySourceUrl(search.SourceURL)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}
	if source.Name != "" {
		records, err = w.storage.ListRecordsSourceID(source.ID)

		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	dataTelegram, err := w.tg.ParseChatTelegram(search.SourceURL, limit, records)
	if err != nil {
		w.log.Error("search", zap.Error(err))
	}

	if source.Name == "" {
		err = w.storage.CreateSource(dataTelegram.Source)
		if err != nil {
			w.log.Error("error: create", zap.Error(err))
			return err
		}
		source, err = w.storage.SearchBySourceUrl(dataTelegram.Source.SourceUrl)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}

	}

	var recordsFilter = make([]models.RecordDTO, 0, len(dataTelegram.Records))
	for _, r := range dataTelegram.Records {
		recordFilter := r

		existText := false
		if len(records) != 0 {
			for _, record := range records { //  ToDo: optimize
				if record.RecordText == r.RecordText {
					r.RecordText = record.RecordText
					existText = true
					break
				}

			}
		}

		if existText {
			continue
		}
		recordsFilter = append(recordsFilter, recordFilter)
	}

	dataTelegram.Records = recordsFilter

	for i := range dataTelegram.Records {
		dataTelegram.Records[i].SourceID = source.ID
	}

	dataTelegram.Source = source
	if len(dataTelegram.Records) != 0 {
		err = w.storage.CreateRecords(dataTelegram.Records)
		if err != nil {
			w.log.Error("error: create", zap.Error(err))
			return err
		}
	}

	records, err = w.storage.ListRecordsSourceID(source.ID)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	dataTelegram.Records = records[:limit]

	switch search.ScoreType {
	case models.Racism:
		err = w.parseSourceTypeRacism(search, dataTelegram, updateChan)

	case models.Waffler:
		err = w.parseSourceTypeWaffler(search, dataTelegram, updateChan)

	}

	return err
}

type wafflerRecords struct {
	sync.Mutex
	records []models.WafflerDTO
}

func (w *WafflerService) parseSourceTypeWaffler(search *message.ParserRequest, dataTelegram *data_source.DataTelegram, updateChan chan<- bool) error {
	g := errgroup.Group{}
	g.SetLimit(20)

	newWafflerRecords := wafflerRecords{}
	lanModel := w.createLanModel(search)
	for _, r := range dataTelegram.Records {
		r := r
		for _, r2 := range dataTelegram.Records {
			if r.RecordText == r2.RecordText || r2.CreatedTs.Before(r.CreatedTs) {
				continue
			}
			r2 := r2

			records, err := w.storage.ListWafflerRecords(&models.WafflerDTO{
				RecordIDBefore: r.ID,
				RecordIDAfter:  r2.ID,
			})
			if err != nil {
				w.log.Error("error: search", zap.Error(err))
				return err
			}
			if len(records) != 0 {
				updateChan <- true
				continue
			}

			text := r.RecordText + " и " + r2.RecordText

			g.Go(func() error {
				var err error
				res, err := lanModel.ConstructQuestionGPT(text, search.ScoreType)
				if res != nil {
					newWafflerRecords.Lock()
					defer newWafflerRecords.Unlock()
					score := models.NewNullInt64(int64(*res))
					if score.Int64 < 0 {
						score = models.NullInt64{}
					}
					newWafflerRecords.records = append(newWafflerRecords.records, models.WafflerDTO{
						Score:           score,
						ParserType:      models.GPT3_5TURBO,
						RecordIDBefore:  r.ID,
						RecordIDAfter:   r.ID,
						CreatedTsBefore: r.CreatedTs,
						CreatedTsAfter:  r2.CreatedTs,
						SourceID:        dataTelegram.Source.ID,
					})
					return nil
				} else {
					if err != nil {
						w.log.Error("error: search", zap.Error(err))
						return err
					}
				}

				return nil
			})

		}

	}
	err := g.Wait()
	if err != nil {
		w.log.Warn("error: search", zap.Error(err))
		return err
	}

	err = w.storage.CreateWafflerRecords(newWafflerRecords.records)
	if err != nil {
		w.log.Error("error: search", zap.Error(err))
		return err
	}

	if len(dataTelegram.Records) != 0 {

		wafflerRecords, err := w.storage.ListWafflerRecords(&models.WafflerDTO{SourceID: dataTelegram.Source.ID})
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}

		dataTelegram.Source.WafflerScore = averageWafflerScore(wafflerRecords)

		err = w.storage.UpdateSource(dataTelegram.Source)
		if err != nil {
			w.log.Error("error: search", zap.Error(err))
			return err
		}
	}

	return err
}

//  ToDo: racizm rename to racism

func (w *WafflerService) Search(search *message.Search) (*message.SearchResponse, error) {
	res := &message.SearchResponse{
		Sources: make([]models.SourceDTO, 0, search.Limit),
	}

	orders, err := convertSourceOrder(search.Order)
	if err != nil {
		return nil, err
	}

	if search.Cursor.Partition == 0 {
		res.Sources, err = w.storage.
			SearchLikeBySourceName(search.QueryForName, search.SourceType, search.Cursor.Offset, orders, search.Limit)

		if err != nil {
			return nil, err
		}

		if len(res.Sources) == search.Limit {
			search.Cursor.Offset += len(res.Sources)
		} else {
			search.Cursor.Partition = 1
			search.Limit -= len(res.Sources)
			search.Cursor.Offset = 0
		}

	}

	if search.Cursor.Partition == 1 {
		resURL, err := w.storage.
			SearchLikeBySourceURLNotName(search.QueryForName, search.SourceType, search.Cursor.Offset, orders, search.Limit)
		if err != nil {
			return nil, err
		}

		res.Sources = append(res.Sources, resURL...)

		if len(resURL) == search.Limit {
			search.Cursor.Offset += len(resURL)
		} else {
			search.Cursor = nil
		}
	}

	res.Cursor = search.Cursor

	return res, nil
}

func (w *WafflerService) createLanModel(search *message.ParserRequest) language_model.LanguageModel {
	var lanModel language_model.LanguageModel
	if search.Parser.Type == models.YakiModel_GPT3_5TURBO {
		lanModel = w.gpt
	} else {
		gptInstance := gpt.NewAiLanguageModel(search.Parser.Token)
		lanModel = language_model.NewChatGPTWrapper(gptInstance, w.log)
	}
	return lanModel
}

func (s *WafflerService) PriceSource(request *message.PriceRequest) (*message.PriceResponse, error) {
	response := message.PriceResponse{Currency: request.Currency}
	var price float64

	switch request.ScoreType {
	case models.Waffler:
		price = priceWaffler(request.Limit)
	case models.Racism:
		price = priceRacism(request.Limit)
	}

	if request.Currency == "RUB" {
		price *= 100
	}

	if message.ValidateParser(int(request.Parser.Type)) {
		price /= 2
	}

	response.Price = fmt.Sprintf("%.3f", price)

	return &response, nil
}
