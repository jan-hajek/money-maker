package app

import "errors"

type Resolver interface {
	Resolve(input ResolverInput) ResolverResult
	GetCalculators() []Calculator
	GetPrintValues() []PrintValue
}

type ResolverInput struct {
	DateInput         DateInput
	CalculatorResults map[string]CalculatorResult
	History           *History
	Position          *Position
}

func (s *ResolverInput) CalculatorResult(name string) CalculatorResult {
	return s.CalculatorResults[name]
}

type ResolverResult struct {
	Action       ResolverAction
	PositionType PositionType
	Amount       float64
	Sl           float64
	Costs        float64
}

type ResolverAction string

const (
	SKIP   ResolverAction = "skip"
	OPEN                  = "open"
	CLOSE                 = "close"
	CHANGE                = "change"
)

type ResolverFactory interface {
	GetName() string
	GetDefaultConfig(config map[string]map[string]interface{}) ResolverFactoryConfig
	GetBatchConfigs(config map[string]map[string]interface{}) []ResolverFactoryConfig
	Create(config ResolverFactoryConfig) Resolver
}

type ResolverFactoryConfig interface {
}

type ResolverFactoryRegistry struct {
	Items map[string]ResolverFactory
}

func (s *ResolverFactoryRegistry) Add(r ResolverFactory) {
	s.Items[r.GetName()] = r
}

func (s *ResolverFactoryRegistry) GetByName(name string) (ResolverFactory, error) {
	item, ok := s.Items[name]
	if ok == false {
		return nil, errors.New("unknown resolver factory " + name)
	}
	return item, nil
}
