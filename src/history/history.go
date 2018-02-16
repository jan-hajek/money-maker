package history

import (
	"errors"
	"github.com/jelito/money-maker/src/interfaces"
)

type History struct {
	Strategy   interfaces.Strategy
	Indicators []interfaces.Indicator
	items      []interfaces.HistoryItem
}

func (s *History) GetStrategy() interfaces.Strategy {
	return s.Strategy
}

func (s *History) AddItem(result interfaces.HistoryItem) {
	s.items = append(s.items, result)
}

// TODO - jhajek error
func (s *History) GetLastItems(numOfLast int) []interfaces.HistoryItem {
	count := len(s.items)
	if numOfLast < count {
		return s.items[count-numOfLast:]
	} else {
		return s.items
	}
}

func (s *History) GetLastItem() (interfaces.HistoryItem, error) {
	items := s.GetLastItems(1)
	if len(items) == 1 {
		return items[0], nil
	} else {
		return &HistoryItem{}, errors.New("last item not found")
	}
}

func (s *History) GetAll() []interfaces.HistoryItem {
	return s.items
}
