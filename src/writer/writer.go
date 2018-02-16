package writer

import (
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/summary"
)

type Writer struct {
	Outputs []Output
}

func (s *Writer) Open() error {
	for _, output := range s.Outputs {
		err := output.Open()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteHistory(historyItems []interfaces.HistoryItem) error {
	for _, output := range s.Outputs {
		err := output.WriteHistory(historyItems)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteSummaryHeader(summary *summary.Summary) error {
	for _, output := range s.Outputs {
		err := output.WriteSummaryHeader(summary)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteSummaryRow(summary *summary.Summary) error {
	for _, output := range s.Outputs {
		err := output.WriteSummaryRow(summary)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) Close() error {
	for _, output := range s.Outputs {
		err := output.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
