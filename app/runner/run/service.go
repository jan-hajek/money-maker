package run

import (
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

		go func() {
			for range ticker.C {
				s.runTitleCron(t, tradesPerTitle)
			}
		}()

	}

	wg.Wait()
}

func (s *Service) runTitleCron(
	t *entity.Title,
	tradesPerTitle map[string][]*appTrade.Service,
) {
	titleLog := s.Log.WithField("title", t.Name)
	titleLog.Info("download price")

	dateInput, err := s.Registry.GetByName(t.ClassName).(interfaces.TitleFactory).Create(t).LoadLast()
	if err != nil {
		titleLog.Error(err)
	}

	exists, err := s.PriceRepository.GetByTitleAndDate(t.Id, dateInput.Date)
	if err != nil {
		titleLog.Error(err)
	}

	if exists != nil {
		titleLog.Info("price not change")
		return
	}

	titleLog.Info("save price to db")
	s.savePriceToDb(t, dateInput)

	if trades, exists := tradesPerTitle[t.Id]; exists {
		// TODO - jhajek routines
		for _, tr := range trades {
			tradeLog := titleLog.WithField("trade", tr.Trade.Id)

			history, err := tr.Run(dateInput)
			if err != nil {
				tradeLog.Error(err)
			}

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

			err = s.Writer.Open()
			if err != nil {
				tradeLog.Error(err)
			}
			s.Writer.WriteHistory(lastHistoryItems)
			err = s.Writer.Close()
			if err != nil {
				tradeLog.Error(err)
			}

			s.Log.Info("----------------")
		}
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
	err := s.PriceRepository.Insert(priceEnt)
	if err != nil {
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
	s.Log.WithField("title", t.Name).Info("load missing prices")

	if s.DownloadMissingPrices {
		dateInputs, err := s.Registry.GetByName(t.ClassName).(interfaces.TitleFactory).Create(t).LoadDataFrom()
		if err != nil {
			s.Log.WithField("title", t.Name).Error(err)
		}

		for _, dateInput := range dateInputs {
			s.savePriceToDb(t, dateInput)
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

	if len(lastPrices) != limit {
		s.Log.
			WithField("title", t.Name).
			WithField("limit", limit).
			Fatal("title needs more prices for warm up")
	} else {
		s.Log.
			WithField("title", t.Name).
			WithField("count", limit).
			Info("warm up title")
	}

	for _, pr := range lastPrices {
		dateInput := dateInput.CreateFromEntity(pr)

		if trades, exists := tradesPerTitle[t.Id]; exists {
			for _, tr := range trades {
				tr.Run(dateInput)
			}
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
