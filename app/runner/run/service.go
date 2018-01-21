package run

import (
	"errors"
	"fmt"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/dateInput"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/interfaces"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/mailer"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/repository/strategy"
	"github.com/jelito/money-maker/app/repository/title"
	"github.com/jelito/money-maker/app/repository/trade"
	appTrade "github.com/jelito/money-maker/app/trade"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Service struct {
	Registry              *registry.Registry
	PriceRepository       *price.Service
	StrategyRepository    *strategy.Service
	TitleRepository       *title.Service
	TradeRepository       *trade.Service
	TradeFactory          *appTrade.Factory
	Log                   log.Log
	Writer                *app.Writer
	Mailer                *mailer.Service
	DownloadMissingPrices bool
}

func (s *Service) Run() {

	wg := sync.WaitGroup{}

	tradesPerTitle := s.getTradesPerTitle()

	for _, t := range s.getTitles() {
		wg.Add(1)

		s.downloadMissingPrices(t)

		s.warmUpTrades(t, tradesPerTitle)

		duration := time.Second * time.Duration(t.DownloadInterval)
		ticker := time.NewTicker(duration)
		s.Log.WithField("title", t.Name).WithField("duration", duration.String()).Info("start watching title")

		trades, exists := tradesPerTitle[t.Id]
		if !exists {
			trades = make([]*appTrade.Service, 0)
		}

		go func() {
			for range ticker.C {
				s.runTitleCron(t, trades)
			}
		}()
	}

	wg.Wait()
}

func (s *Service) runTitleCron(t *entity.Title, trades []*appTrade.Service) {
	titleLog := s.Log.WithField("title", t.Name)
	titleLog.Info("download price")

	dateInput, err := s.Registry.GetByName(t.ClassName).(interfaces.TitleFactory).Create(t).LoadLast()
	if err != nil {
		titleLog.Error(err)
		return
	}

	storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dateInput.Date)
	if err != nil {
		titleLog.Error(err)
		return
	}

	if storedPrice != nil {
		if err = s.checkStoredPrice(storedPrice, dateInput); err != nil {
			titleLog.Warning(err)
		} else {
			titleLog.Info("price not change")
		}
		return
	}

	titleLog.Info("save price to db")
	s.savePriceToDb(t, dateInput)

	// TODO - jhajek routines
	for _, tr := range trades {
		tradeLog := titleLog.WithField("trade", tr.Trade.Id)

		history, err := tr.Run(dateInput)
		if err != nil {
			tradeLog.Error(err)
			continue
		}
		s.writeLastHistoryItems(history, tradeLog)

		lastHistoryItem, err := history.GetLastItem()
		if err != nil {
			tradeLog.Error(err)
			return
		}
		s.Mailer.Send(
			fmt.Sprintf("title: %s, action: %s", t.Name, lastHistoryItem.StrategyResult.Action),
			fmt.Sprintf("title: %s, action: %s", t.Name, lastHistoryItem.StrategyResult.Action),
		)

		s.Log.Info("----------------")
	}
}

func (s *Service) savePriceToDb(t *entity.Title, dateInput app.DateInput) {
	priceEnt := &entity.Price{
		Id:         uuid.Must(uuid.NewV4()).String(),
		TitleId:    t.Id,
		Date:       dateInput.Date,
		OpenPrice:  dateInput.OpenPrice.Val(),
		HighPrice:  dateInput.HighPrice.Val(),
		LowPrice:   dateInput.LowPrice.Val(),
		ClosePrice: dateInput.ClosePrice.Val(),
	}
	if err := s.PriceRepository.Insert(priceEnt); err != nil {
		s.Log.Error(err)
	}
}

func (s *Service) getTitles() []*entity.Title {
	list, err := s.TitleRepository.GetAllActive()
	if err != nil {
		s.Log.Fatal(err)
	}

	return list
}

func (s *Service) downloadMissingPrices(t *entity.Title) {
	if s.DownloadMissingPrices {
		log := s.Log.WithField("title", t.Name)
		log.Info("load missing prices")

		dateInputs, err := s.Registry.GetByName(t.ClassName).(interfaces.TitleFactory).Create(t).LoadDataFrom()
		if err != nil {
			log.Fatal(err)
		}

		for _, dateInput := range dateInputs {
			storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dateInput.Date)
			if err != nil {
				log.Fatal(err)
			}
			if storedPrice == nil {
				s.savePriceToDb(t, dateInput)
			} else {
				if err = s.checkStoredPrice(storedPrice, dateInput); err != nil {
					log.Warning(err)
				}
			}
		}
	}
}

func (s *Service) warmUpTrades(
	t *entity.Title,
	tradesPerTitle map[string][]*appTrade.Service,
) {
	// TODO - jhajek config
	limit := 100
	lastPrices := s.getLastPrices(t.Id, limit)

	log := s.Log.WithField("title", t.Name)

	if len(lastPrices) != limit {
		log.WithField("limit", limit).Fatal("title needs more prices for warm up")
	} else {
		log.WithField("count", limit).Info("warm up title")
	}

	if trades, exists := tradesPerTitle[t.Id]; exists {
		for _, tr := range trades {
			var history *app.History
			var err error

			for _, pr := range lastPrices {
				dateInput := dateInput.CreateFromEntity(pr)

				history, err = tr.Run(dateInput)
				if err != nil {
					log.Fatal(err)
				}
			}
			s.writeLastHistoryItems(history, log)
		}
	}
}

func (s *Service) getLastPrices(titleId string, limit int) []*entity.Price {
	list, err := s.PriceRepository.GetLastItemsByTitle(titleId, limit)
	if err != nil {
		s.Log.Fatal(err)
	}

	return list
}

func (s *Service) getTradesPerTitle() map[string][]*appTrade.Service {

	list, err := s.TradeRepository.GetAllActive()
	if err != nil {
		s.Log.Fatal(err)
	}

	trades := make(map[string][]*appTrade.Service)
	for _, t := range list {
		strategyEntity := s.getStrategy(t.StrategyId)
		strategyFactory := s.Registry.GetByName(strategyEntity.ClassName).(app.StrategyFactory)

		strategyClass := strategyFactory.Create(
			strategyFactory.GetDefaultConfig(t.Params.Data),
		)

		if _, exists := trades[t.TitleId]; exists == false {
			trades[t.TitleId] = make([]*appTrade.Service, 0)
		}

		trades[t.TitleId] = append(
			trades[t.TitleId],
			s.TradeFactory.Create(t, strategyClass),
		)
	}

	return trades
}

func (s *Service) getStrategy(id string) *entity.Strategy {
	ent, err := s.StrategyRepository.GetById(id)
	if err != nil {
		s.Log.Fatal(err)
	}

	return ent
}

func (s *Service) writeLastHistoryItems(history *app.History, tradeLog *logrus.Entry) {
	lastHistoryItems := history.GetLastItems(2)

	if lastHistoryItems[0].Position == nil {
		tradeLog.Info("no last position")
	} else {
		tradeLog.
			WithField("lastPosition", lastHistoryItems[0].Position.Id).
			Info("use last position")
	}

	tradeLog.
		WithField("action", lastHistoryItems[1].StrategyResult.Action).
		WithField("type", lastHistoryItems[1].StrategyResult.PositionType).
		Info("strategy result")

	if err := s.Writer.Open(); err != nil {
		tradeLog.Error(err)
	}
	if err := s.Writer.WriteHistory(lastHistoryItems); err != nil {
		tradeLog.Error(err)
	}
	if err := s.Writer.Close(); err != nil {
		tradeLog.Error(err)
	}
}

func (s *Service) checkStoredPrice(storedPrice *entity.Price, dateInput app.DateInput) error {
	type diff struct {
		field       string
		storedPrice float64
		newPrice    float64
	}

	diffList := make([]diff, 0)

	checkField := func(field string, storedPrice, newPrice float64) {
		if storedPrice != newPrice {
			diffList = append(diffList, diff{
				field:       field,
				storedPrice: storedPrice,
				newPrice:    newPrice,
			})
		}
	}

	checkField("open", storedPrice.OpenPrice, dateInput.OpenPrice.Val())
	checkField("high", storedPrice.HighPrice, dateInput.HighPrice.Val())
	checkField("low", storedPrice.LowPrice, dateInput.LowPrice.Val())
	checkField("close", storedPrice.ClosePrice, dateInput.ClosePrice.Val())

	if len(diffList) > 0 {
		message := fmt.Sprintf("stored price is different than new date: %s", dateInput.Date.Format("2006-01-02 15:04 -0700"))
		for _, diff := range diffList {
			message += fmt.Sprintf("field: %s, stored price: %.3f, new price: %.3f", diff.field, diff.storedPrice, diff.newPrice)
		}
		return errors.New(message)

	}

	return nil
}
