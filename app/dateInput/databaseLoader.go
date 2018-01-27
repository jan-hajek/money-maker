package dateInput

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/repository/price"
)

type DatabaseLoader struct {
	TitleId         string
	PriceRepository *price.Service
}

func (s *DatabaseLoader) Load() ([]app.DateInput, error) {
	prices, err := s.PriceRepository.GetLastItemsByTitle(s.TitleId, 10000)
	if err != nil {
		return nil, err
	}
	dateInputs := make([]app.DateInput, len(prices))
	for index, pr := range prices {
		dateInputs[index] = CreateFromEntity(pr)
	}

	return dateInputs, nil
}
