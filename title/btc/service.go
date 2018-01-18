package btc

import (
	"encoding/json"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"io/ioutil"
	"log"
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
}

func (s *Service) LoadLast() app.DateInput {
	jsonLines := s.load()
	lastResult := jsonLines[len(jsonLines)-1]

	return app.DateInput{
		Date:       time.Unix(lastResult.Time, 0),
		OpenPrice:  float.New(lastResult.Open),
		HighPrice:  float.New(lastResult.High),
		LowPrice:   float.New(lastResult.Low),
		ClosePrice: float.New(lastResult.Close),
	}
}

func (s *Service) LoadDataFrom() []app.DateInput {
	jsonLines := s.load()

	a := make([]app.DateInput, 0)

	for _, line := range jsonLines {
		a = append(a, app.DateInput{
			Date:       time.Unix(line.Time, 0),
			OpenPrice:  float.New(line.Open),
			HighPrice:  float.New(line.High),
			LowPrice:   float.New(line.Low),
			ClosePrice: float.New(line.Close),
		})
	}

	return a
}

func (s *Service) load() []JsonLine {
	url := "https://admiralmarkets.com/api/ajax/ticks_cfd?name=BTCUSD&range=day&period=60"

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonLines := make([]JsonLine, 0)
	jsonErr := json.Unmarshal(body, &jsonLines)
	if jsonErr != nil {
		log.Fatal(jsonErr, "---"+string(body)+"---")
	}

	return jsonLines
}
