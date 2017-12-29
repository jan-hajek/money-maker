package app

import (
	"errors"
	"sort"
)

type History struct {
	strategy   Strategy
	indicators []Indicator
	items      []*HistoryItem
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
	DateInput        DateInput
	IndicatorResults map[string]IndicatorResult
	StrategyResult   StrategyResult
	Position         *Position
}

func (s *HistoryItem) IndicatorResult(c Indicator) IndicatorResult {
	item := s.IndicatorResults[c.GetName()]

	return item
}

func (s *HistoryItem) OrderedIndicatorResults() []IndicatorResult {

	var keys []string
	for k := range s.IndicatorResults {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var ordered []IndicatorResult
	for _, k := range keys {
		ordered = append(ordered, s.IndicatorResults[k])
	}

	return ordered
}
