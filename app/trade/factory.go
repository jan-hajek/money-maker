package trade

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/log"
	"github.com/jelito/money-maker/app/position"
)

type Factory struct {
	PositionMaker *position.Maker
	Log           log.Log
}

func (s *Factory) Create(trade *entity.Trade, strategy app.Strategy) *Service {
	indicators := strategy.GetIndicators()

	history := &app.History{
		Strategy:   strategy,
		Indicators: indicators,
	}

	return &Service{
		Trade:         trade,
		strategy:      strategy,
		history:       history,
		indicators:    indicators,
		iteration:     0,
		log:           s.Log,
		positionMaker: s.PositionMaker,
	}
}
