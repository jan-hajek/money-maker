package run

import (
	"errors"
	"fmt"

	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/title"
	"github.com/satori/go.uuid"
)

func (s *Service) savePriceToDb(t *entity.Title, dateInput dateInput.DateInput) {
	priceEnt := &entity.Price{
		Id:         uuid.Must(uuid.NewV4()).String(),
		TitleId:    t.Id,
		Date:       dateInput.Date,
		OpenPrice:  dateInput.OpenPrice.Val(),
		HighPrice:  dateInput.HighPrice.Val(),
		LowPrice:   dateInput.LowPrice.Val(),
		ClosePrice: dateInput.ClosePrice.Val(),
	}
	if err := s.PriceRepository.Insert(priceEnt); err != nil {
		s.Log.Error(err)
	}
}

func (s *Service) downloadMissingPrices(t *entity.Title) {
	if s.DownloadMissingPrices {
		titleLog := s.Log.WithField("title", t.Name)
		titleLog.Info("load missing prices")

		dateInputs, err := s.Registry.GetByName(t.ClassName).(title.Factory).Create(t).LoadDataFrom()
		if err != nil {
			titleLog.Fatal(err)
		}

		for _, dtInput := range dateInputs {
			storedPrice, err := s.PriceRepository.GetByTitleAndDate(t.Id, dtInput.Date)
			if err != nil {
				titleLog.Fatal(err)
			}
			if storedPrice == nil {
				s.savePriceToDb(t, dtInput)
			} else {
				if err = s.checkStoredPrice(storedPrice, dtInput); err != nil {
					titleLog.Warning(err)
				}
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

func (s *Service) checkStoredPrice(storedPrice *entity.Price, dateInput dateInput.DateInput) error {
	type diff struct {
		field       string
		storedPrice float64
		newPrice    float64
	}

	diffList := make([]diff, 0)

	checkField := func(field string, storedPrice, newPrice float64) {
		if storedPrice != newPrice {
			diffList = append(diffList, diff{
				field:       field,
				storedPrice: storedPrice,
				newPrice:    newPrice,
			})
		}
	}

	checkField("open", storedPrice.OpenPrice, dateInput.OpenPrice.Val())
	checkField("high", storedPrice.HighPrice, dateInput.HighPrice.Val())
	checkField("low", storedPrice.LowPrice, dateInput.LowPrice.Val())
	checkField("close", storedPrice.ClosePrice, dateInput.ClosePrice.Val())

	if len(diffList) > 0 {
		message := fmt.Sprintf("price diff [stored <> new], date: %s", dateInput.Date.Format("2006-01-02 15:04 -0700"))
		for _, diff := range diffList {
			message += fmt.Sprintf(", %s [%.3f <> %.3f]", diff.field, diff.storedPrice, diff.newPrice)
		}
		return errors.New(message)

	}

	return nil
}
