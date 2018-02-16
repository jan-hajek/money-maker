package admiralMarkets

import (
	"encoding/json"
	"errors"
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/math/float"
	"io/ioutil"
	"net/http"
	"time"
)

// [{"time":1515193200,"open":16552.3,"high":16875.1,"low":16551.8,"volume":4343,"close":16770.8,"date":"05 Jan, 2018 23:00 EET"}]

type JsonLine struct {
	Time  int64
	Open  float64
	High  float64
	Low   float64
	Close float64
	Date  string
}

type Service struct {
	titleEntity *entity.Title
}

func (s *Service) LoadLast() (dateInput.DateInput, error) {
	jsonLines, err := s.load(s.titleEntity.DataUrl)
	if err != nil {
		return dateInput.DateInput{}, err
	}
	lastResult := jsonLines[len(jsonLines)-1]

	return dateInput.DateInput{
		Date:       time.Unix(lastResult.Time, 0),
		OpenPrice:  float.New(lastResult.Open),
		HighPrice:  float.New(lastResult.High),
		LowPrice:   float.New(lastResult.Low),
		ClosePrice: float.New(lastResult.Close),
	}, nil
}

func (s *Service) LoadDataFrom() ([]dateInput.DateInput, error) {
	jsonLines, err := s.load(s.titleEntity.BatchDataUrl)
	if err != nil {
		return nil, err
	}
	a := make([]dateInput.DateInput, 0)

	for _, line := range jsonLines {
		a = append(a, dateInput.DateInput{
			Date:       time.Unix(line.Time, 0),
			OpenPrice:  float.New(line.Open),
			HighPrice:  float.New(line.High),
			LowPrice:   float.New(line.Low),
			ClosePrice: float.New(line.Close),
		})
	}

	return a, nil
}

func (s *Service) load(url string) ([]JsonLine, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	jsonLines := make([]JsonLine, 0)
	jsonErr := json.Unmarshal(body, &jsonLines)
	if jsonErr != nil {
		return nil, errors.New(string(body))
	}

	return jsonLines, nil
}
