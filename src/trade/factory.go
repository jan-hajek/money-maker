package trade

import (
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/history"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/log"
	"github.com/jelito/money-maker/src/position"
)

type Factory struct {
	PositionMaker *position.PositionMaker
	Log           log.Log
}

func (s *Factory) Create(trade *entity.Trade, strategy interfaces.Strategy) *Service {
	indicators := strategy.GetIndicators()

	history := &history.History{
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
		idGenerator:   &position.UuidGenerator{},
	}
}
