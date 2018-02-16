package plus500

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/float"
	"io/ioutil"
	"net/http"
	"time"
)

// [1516935600000,11561.44,11561.94,11420.44,11480.44]

type RawLine []interface{}

type JsonLine struct {
	Time  int64
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type Service struct {
	titleEntity *entity.Title
}

func (s *Service) LoadLast() (app.DateInput, error) {
	jsonLines, err := s.load(s.titleEntity.DataUrl)
	if err != nil {
		return app.DateInput{}, err
	}
	lastResult := jsonLines[len(jsonLines)-1]

	return app.DateInput{
		Date:       time.Unix(lastResult.Time, 0).In(time.Local),
		OpenPrice:  float.New(lastResult.Open),
		HighPrice:  float.New(lastResult.High),
		LowPrice:   float.New(lastResult.Low),
		ClosePrice: float.New(lastResult.Close),
	}, nil
}

func (s *Service) LoadDataFrom() ([]app.DateInput, error) {
	jsonLines, err := s.load(s.titleEntity.BatchDataUrl)
	if err != nil {
		return nil, err
	}
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

	// remove " from start & end, because input is json in xml string :/
	body = body[1 : len(body)-1]

	rawLines := make([]RawLine, 0)
	jsonErr := json.Unmarshal(body, &rawLines)
	if jsonErr != nil {
		return nil, errors.New(string(body))
	}

	jsonLines := make([]JsonLine, 0)
	for _, line := range rawLines {
		fmt.Println(time.Unix(int64(line[0].(float64)/1000), 0).In(time.Local))

		jsonLines = append(jsonLines, JsonLine{
			Time:  int64(line[0].(float64) / 1000),
			Open:  line[1].(float64),
			High:  line[2].(float64),
			Low:   line[3].(float64),
			Close: line[4].(float64),
		})
	}

	return jsonLines, nil
}
