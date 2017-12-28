package app

import (
	"errors"
	"sort"
)

type History struct {
	resolver    Resolver
	calculators []Calculator
	items       []*HistoryItem
}

func (s *History) AddItem(result *HistoryItem) {
	s.items = append(s.items, result)
}

func (s *History) GetLastItems(numOfLast int) []*HistoryItem {
	count := len(s.items)
	if numOfLast < count {
		return s.items[count-numOfLast : count]
	} else {
		return s.items
	}
}

func (s *History) GetLastItem() (*HistoryItem, error) {
	items := s.GetLastItems(1)
	if len(items) == 1 {
		return items[0], nil
	} else {
		return &HistoryItem{}, errors.New("last item not found")
	}
}

func (s *History) GetAll() []*HistoryItem {
	return s.items
}

type HistoryItem struct {
	DateInput         DateInput
	CalculatorResults map[string]CalculatorResult
	ResolverResult    ResolverResult
	Position          *Position
}

func (s *HistoryItem) CalculatorResult(name string) CalculatorResult {
	item := s.CalculatorResults[name]

	return item
}

func (s *HistoryItem) OrderedCalculatorResults() []CalculatorResult {

	var keys []string
	for k := range s.CalculatorResults {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var ordered []CalculatorResult
	for _, k := range keys {
		ordered = append(ordered, s.CalculatorResults[k])
	}

	return ordered
}
