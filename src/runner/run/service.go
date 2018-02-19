package run

import (
	"errors"
	"fmt"
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/log"
	"github.com/jelito/money-maker/src/mailer"
	"github.com/jelito/money-maker/src/registry"
	"github.com/jelito/money-maker/src/repository/price"
	strategyRepo "github.com/jelito/money-maker/src/repository/strategy"
	titleRepo "github.com/jelito/money-maker/src/repository/title"
	"github.com/jelito/money-maker/src/repository/trade"
	"github.com/jelito/money-maker/src/strategy"
	"github.com/jelito/money-maker/src/title"
	appTrade "github.com/jelito/money-maker/src/trade"
	"github.com/jelito/money-maker/src/writer"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Service struct {
	Registry              *registry.Registry
	PriceRepository       *price.Service
	StrategyRepository    *strategyRepo.Service
	TitleRepository       *titleRepo.Service
	TradeRepository       *trade.Service
	TradeFactory          *appTrade.Factory
	Log                   log.Log
	Writer                *writer.Writer
	MailBufferFactory     *mailer.BufferFactory
	DownloadMissingPrices bool
}

func (s *Service) Run() {
	wg := sync.WaitGroup{}
	tradesPerTitle, tradesCount := s.getTradesPerTitle()
	titles := s.getTitles()
	mailBuffer := s.MailBufferFactory.Create(tradesCount)

	for _, t := range titles {
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

		go func(title2 *entity.Title, trades2 []*appTrade.Service) {
			for range ticker.C {
				s.runTitleCron(title2, trades2, mailBuffer)
			}
		}(t, trades)
	}

	wg.Wait()
}

func (s *Service) runTitleCron(
	t *entity.Title,
	trades []*appTrade.Service,
	mailBuffer chan<- mailer.BufferItem,
) {
	titleLog := s.Log.WithField("title", t.Name)
	titleLog.Debug("download price")

	dtInput, err := s.Registry.GetByName(t.ClassName).(title.Factory).Create(t).LoadLast()
	if err != nil {
		titleLog.Error(err)
		return
	}

	storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dtInput.Date)
	if err != nil {
		titleLog.Error(err)
		return
	}

	if storedPrice != nil {
		if err = s.checkStoredPrice(storedPrice, dtInput); err != nil {
			titleLog.Warning(err)
		} else {
			titleLog.Debug("no new price")
		}
		return
	}

	titleLog.Debug("save price to db")
	s.savePriceToDb(t, dtInput)

	// TODO - jhajek routines
	for _, tr := range trades {
		tradeLog := titleLog.WithField("trade", tr.Trade.Id)

		history, err := tr.Run(dtInput)
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

		s.sendEmail(lastHistoryItem, mailBuffer, t)

		s.Log.Info("----------------")
	}
}

func (s *Service) savePriceToDb(t *entity.Title, dateInput dateInput.DateInput) {
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
		titleLog := s.Log.WithField("title", t.Name)
		titleLog.Info("load missing prices")

		dateInputs, err := s.Registry.GetByName(t.ClassName).(title.Factory).Create(t).LoadDataFrom()
		if err != nil {
			titleLog.Fatal(err)
		}

		for _, dtInput := range dateInputs {
			storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dtInput.Date)
			if err != nil {
				titleLog.Fatal(err)
			}
			if storedPrice == nil {
				s.savePriceToDb(t, dtInput)
			} else {
				if err = s.checkStoredPrice(storedPrice, dtInput); err != nil {
					titleLog.Warning(err)
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

	titleLog := s.Log.WithField("title", t.Name)

	if len(lastPrices) != limit {
		titleLog.WithField("limit", limit).Fatal("title needs more prices for warm up")
	} else {
		titleLog.WithField("count", limit).Info("warm up title")
	}

	if trades, exists := tradesPerTitle[t.Id]; exists {
		for _, tr := range trades {
			var history interfaces.History
			var err error

			for _, pr := range lastPrices {
				history, err = tr.Run(dateInput.CreateFromEntity(pr))
				if err != nil {
					titleLog.Fatal(err)
				}
			}
			s.writeLastHistoryItems(history, titleLog)
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

func (s *Service) getTradesPerTitle() (map[string][]*appTrade.Service, int) {

	list, err := s.TradeRepository.GetAllActive()
	if err != nil {
		s.Log.Fatal(err)
	}

	trades := make(map[string][]*appTrade.Service)
	count := 0
	for _, t := range list {
		strategyEntity := s.getStrategy(t.StrategyId)
		strategyFactory := s.Registry.GetByName(strategyEntity.ClassName).(strategy.Factory)

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
		count++
	}

	return trades, count
}

func (s *Service) getStrategy(id string) *entity.Strategy {
	ent, err := s.StrategyRepository.GetById(id)
	if err != nil {
		s.Log.Fatal(err)
	}

	return ent
}

func (s *Service) writeLastHistoryItems(history interfaces.History, tradeLog *logrus.Entry) {
	lastHistoryItems := history.GetLastItems(2)

	if lastHistoryItems[0].GetPosition() == nil {
		tradeLog.Info("no last position")
	} else {
		tradeLog.
			WithField("lastPosition", lastHistoryItems[0].GetPosition().Id).
			Info("use last position")
	}

	tradeLog.
		WithField("action", lastHistoryItems[1].GetStrategyResult().GetAction()).
		WithField("type", lastHistoryItems[1].GetStrategyResult().GetPositionType()).
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

func (s *Service) checkStoredPrice(storedPrice *entity.Price, dateInput dateInput.DateInput) error {
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
		message := fmt.Sprintf("price diff [stored <> new], date: %s", dateInput.Date.Format("2006-01-02 15:04 -0700"))
		for _, diff := range diffList {
			message += fmt.Sprintf(", %s [%.3f <> %.3f]", diff.field, diff.storedPrice, diff.newPrice)
		}
		return errors.New(message)

	}

	return nil
}

func (s *Service) sendEmail(lastHistoryItem interfaces.HistoryItem, mailBuffer chan<- mailer.BufferItem, t *entity.Title) {
	action := lastHistoryItem.GetStrategyResult().GetAction()
	subject := fmt.Sprintf("%s", action)
	var message string

	if lastHistoryItem.GetStrategyResult().GetReportMessage() != "" {
		message = lastHistoryItem.GetStrategyResult().GetReportMessage()
	} else {
		if lastHistoryItem.GetPosition() == nil {
			message = fmt.Sprintf("%s", action)
		} else {
			message = fmt.Sprintf(
				"%s, %s, %.3f",
				action,
				lastHistoryItem.GetPosition().Type,
				lastHistoryItem.GetPosition().Sl,
			)
		}
	}

	mailBuffer <- mailer.BufferItem{
		Subject: subject,
		Message: fmt.Sprintf("%s, ", t.Name) + message,
	}
}
