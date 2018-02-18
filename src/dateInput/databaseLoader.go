package dateInput

import (
	"github.com/jelito/money-maker/src/repository/price"
)

type DatabaseLoader struct {
	TitleId         string
	PriceRepository *price.Service
}

func (s *DatabaseLoader) Load() ([]DateInput, error) {
	prices, err := s.PriceRepository.GetLastItemsByTitle(s.TitleId, 10000)
	if err != nil {
		return nil, err
	}
	dateInputs := make([]DateInput, len(prices))
	for index, pr := range prices {
		dateInputs[index] = CreateFromEntity(pr)
	}

	return dateInputs, nil
}
