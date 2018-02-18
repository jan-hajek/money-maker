package interfaces

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/printValue"
	"time"
)

type History interface {
	AddItem(result HistoryItem)
	GetLastItem() (HistoryItem, error)
	GetLastItems(numOfLast int) []HistoryItem
	GetStrategy() Strategy
	GetAll() []HistoryItem
}

type HistoryItem interface {
	GetDateInput() dateInput.DateInput
	IndicatorResult(c Indicator) IndicatorResult
	GetStrategyResult() StrategyResult
	GetPosition() *Position
	GetIndicatorResults() map[string]IndicatorResult
}

type Indicator interface {
	Calculate(input IndicatorInput, history History) IndicatorResult
	GetName() string
}

type IndicatorInput struct {
	Date       time.Time
	OpenPrice  float.Float
	ClosePrice float.Float
	HighPrice  float.Float
	LowPrice   float.Float
	Iteration  int
}

type IndicatorResult interface {
	Print() []printValue.PrintValue
}

type Strategy interface {
	Resolve(input StrategyInput) StrategyResult
	GetIndicators() []Indicator
	GetPrintValues() []printValue.PrintValue
}

type StrategyInput interface {
	GetHistory() History
	IndicatorResult(c Indicator) IndicatorResult
	GetDateInput() dateInput.DateInput
	GetPosition() *Position
}

type StrategyResult interface {
	GetAction() StrategyAction
	GetPositionType() PositionType
	GetAmount() float.Float
	GetSl() float.Float
	GetCosts() float.Float
	GetReportMessage() string
}

type StrategyAction string

const (
	SKIP   StrategyAction = "skip"
	OPEN                  = "open"
	CLOSE                 = "close"
	CHANGE                = "change"
)

type Position struct {
	Id                    string
	Type                  PositionType
	StartDate             time.Time
	CloseDate             time.Time
	OpenPrice             float.Float
	ClosePrice            float.Float
	Amount                float.Float
	Sl                    float.Float
	Costs                 float.Float
	Profit                float.Float
	PossibleProfit        float.Float
	PossibleProfitPercent float.Float
}

func (s *Position) Clone() *Position {
	return &Position{
		Id:                    s.Id,
		Type:                  s.Type,
		StartDate:             s.StartDate,
		CloseDate:             s.CloseDate,
		OpenPrice:             s.OpenPrice,
		ClosePrice:            s.ClosePrice,
		Amount:                s.Amount,
		Sl:                    s.Sl,
		Costs:                 s.Costs,
		Profit:                s.Profit,
		PossibleProfit:        s.PossibleProfit,
		PossibleProfitPercent: s.PossibleProfitPercent,
	}
}

type PositionType string

const (
	LONG  PositionType = "long"
	SHORT              = "short"
)
