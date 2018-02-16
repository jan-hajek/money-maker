package dateInput

import (
	"encoding/csv"
	"github.com/jelito/money-maker/src/math/float"
	"io"
	"os"
	"strconv"
	"time"
)

type CsvLoader struct {
	InputFilePath   string
	TimeParseFormat string
}

func (s *CsvLoader) Load() ([]DateInput, error) {

	file, err := os.Open(s.InputFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'

	var data []DateInput

	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		date, err := time.Parse(s.TimeParseFormat, record[0])
		if err != nil {
			return nil, err
		}

		openPrice, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return nil, err
		}

		highPrice, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, err
		}

		lowPrice, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			return nil, err
		}

		closePrice, err := strconv.ParseFloat(record[4], 32)
		if err != nil {
			return nil, err
		}

		data = append(data, DateInput{
			Date:       date,
			OpenPrice:  float.New(openPrice),
			ClosePrice: float.New(closePrice),
			HighPrice:  float.New(highPrice),
			LowPrice:   float.New(lowPrice),
		})

		lineCount += 1
	}

	return data, nil
}
