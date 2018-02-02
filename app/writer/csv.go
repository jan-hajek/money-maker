package writer

import (
	"encoding/csv"
	"github.com/jelito/money-maker/app"
	"os"
)

type CsvWriterOutput struct {
	File       string
	DateFormat string
	file       *os.File
	writer     *csv.Writer
}

func (s *CsvWriterOutput) Open() error {
	var err error
	s.file, err = os.OpenFile(s.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	s.writer = csv.NewWriter(s.file)
	s.writer.Comma = ';'

	return err
}

func (s *CsvWriterOutput) WriteHistory(historyItems []*app.HistoryItem) error {
	if len(historyItems) == 0 {
		return nil
	}
	err := s.writer.Write(writerGetHistoryHeader(historyItems[0]))
	if err != nil {
		return err
	}
	for _, item := range historyItems {
		err = s.writer.Write(writerGetHistoryRow(item, s.DateFormat))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CsvWriterOutput) WriteSummaryHeader(summary *app.Summary) error {
	return s.writer.Write(writerGetSummaryHeader(summary))
}

func (s *CsvWriterOutput) WriteSummaryRow(summary *app.Summary) error {
	return s.writer.Write(writerGetSummaryRow(summary))
}

func (s *CsvWriterOutput) Close() error {
	s.writer.Flush()
	return s.file.Close()
}
