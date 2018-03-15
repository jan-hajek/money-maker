package run

import (
	"fmt"
	"sync"
	"time"

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
	"github.com/sirupsen/logrus"
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

		trades, exists := tradesPerTitle[t.Id]
		if !exists {
			trades = make([]*appTrade.Service, 0)
		}

		go s.runTitleWaiting(t, trades, mailBuffer)
	}

	wg.Wait()
}

func (s *Service) runTitleWaiting(
	t *entity.Title,
	trades []*appTrade.Service,
	mailBuffer chan<- mailer.BufferItem,
) {
	titleLog := s.Log.WithField("title", t.Name)

	for {
		now := time.Now()
		// TODO - jhajek variabilni
		nextHourStart := now.Add(time.Hour).Truncate(time.Hour)
		nextTick := nextHourStart.Sub(now)

		titleLog.Info("next processing in ", nextTick)

		<-time.NewTimer(nextTick).C

		s.runTitlePriceDownloading(t, titleLog, trades, mailBuffer)
	}
}

func (s *Service) runTitlePriceDownloading(
	t *entity.Title,
	log log.Log,
	trades []*appTrade.Service,
	mailBuffer chan<- mailer.BufferItem,
) {
	tickerD := time.Minute * 1
	timeoutD := time.Minute * 10

	ticker := time.NewTicker(tickerD)
	timeout := time.NewTimer(timeoutD)

	defer ticker.Stop()
	defer timeout.Stop()

	for {
		log.Info("downloading price")

		dtInput, err := s.getNewTitlePrice(t)
		if err != nil {
			log.Error(err)
			return
		}

		if s.isTitlePriceNew(t, dtInput, log) {
			log.Info("saving price to db")
			s.savePriceToDb(t, dtInput)

			s.processTitle(t, dtInput, trades, log, mailBuffer)
			return
		}

		select {
		case <-timeout.C:
			log.Warning("timeout, no strategy ")
			return
		case <-ticker.C:
		}
	}
}
func (s *Service) getNewTitlePrice(t *entity.Title) (dateInput.DateInput, error) {
	return s.Registry.GetByName(t.ClassName).(title.Factory).Create(t).LoadLast()
}

func (s *Service) isTitlePriceNew(t *entity.Title, dtInput dateInput.DateInput, log log.Log) bool {
	storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dtInput.Date)
	if err != nil {
		log.Error(err)
		return false
	}

	if storedPrice != nil {
		if err = s.checkStoredPrice(storedPrice, dtInput); err != nil {
			log.Warning(err)
		} else {
			log.Info("no new price")
		}
		return false
	}

	return true
}

func (s *Service) processTitle(
	t *entity.Title,
	dtInput dateInput.DateInput,
	trades []*appTrade.Service,
	log log.Log,
	mailBuffer chan<- mailer.BufferItem,
) {
	// TODO - jhajek routines
	for _, tr := range trades {
		tradeLog := log.WithField("trade", tr.Trade.Id)

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

func (s *Service) getTitles() []*entity.Title {
	list, err := s.TitleRepository.GetAllActive()
	if err != nil {
		s.Log.Fatal(err)
	}

	return list
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
