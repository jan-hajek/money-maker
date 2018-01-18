package app

import (
	"encoding/csv"
	"fmt"
	"github.com/jelito/money-maker/app/float"
	"io"
	"os"
	"text/tabwriter"
)

type Writer struct {
	Outputs []WriterOutput
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

func (s *Writer) WriteHistory(historyItems []*HistoryItem) error {
	for _, output := range s.Outputs {
		err := output.WriteHistory(historyItems)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteSummaryHeader(summary *Summary) error {
	for _, output := range s.Outputs {
		err := output.WriteSummaryHeader(summary)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteSummaryRow(summary *Summary) error {
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

type WriterOutput interface {
	Open() error
	WriteHistory(historyItems []*HistoryItem) error
	WriteSummaryHeader(summary *Summary) error
	WriteSummaryRow(summary *Summary) error
	Close() error
}

type StdOutWriterOutput struct {
	DateFormat string
	w          *tabwriter.Writer
}

func (s *StdOutWriterOutput) Open() error {
	s.w = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)

	return nil
}

func (s *StdOutWriterOutput) WriteHistory(historyItems []*HistoryItem) error {
	if len(historyItems) == 0 {
		return nil
	}
	s.write(s.w, writerGetHistoryHeader(historyItems[0])...)
	for _, item := range historyItems {
		s.write(s.w, writerGetHistoryRow(item, s.DateFormat)...)
	}
	return nil
}

func (s *StdOutWriterOutput) WriteSummaryHeader(summary *Summary) error {
	s.write(s.w, writerGetSummaryHeader(summary)...)

	return nil
}

func (s *StdOutWriterOutput) WriteSummaryRow(summary *Summary) error {
	s.write(s.w, writerGetSummaryRow(summary)...)
	return nil
}

func (s *StdOutWriterOutput) Close() error {
	return s.w.Flush()
}

func (s *StdOutWriterOutput) write(w io.Writer, a ...string) {
	for _, item := range a {
		fmt.Fprint(w, item, "\t")
	}
	fmt.Fprint(w, "\n")

}

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

func (s *CsvWriterOutput) WriteHistory(historyItems []*HistoryItem) error {
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

func (s *CsvWriterOutput) WriteSummaryHeader(summary *Summary) error {
	return s.writer.Write(writerGetSummaryHeader(summary))
}

func (s *CsvWriterOutput) WriteSummaryRow(summary *Summary) error {
	return s.writer.Write(writerGetSummaryRow(summary))
}

func (s *CsvWriterOutput) Close() error {
	s.writer.Flush()
	return s.file.Close()
}

func writerGetHistoryHeader(item *HistoryItem) []string {
	a := []string{
		"date",
		"price",
	}

	for _, indicatorResult := range item.OrderedIndicatorResults() {
		for _, param := range indicatorResult.Print() {
			a = append(a, param.Label)
		}
	}

	a = append(a,
		"type",
		"id",
		"type",
		"amount",
		"sl",
		"costs",
		"profit",
		"poss. profit",
	)

	return a
}

func writerGetHistoryRow(item *HistoryItem, dateFormat string) []string {
	position := item.Position
	values := []string{
		item.DateInput.Date.Format(dateFormat),
		formatValue(item.DateInput.ClosePrice),
	}

	for _, indicatorResult := range item.OrderedIndicatorResults() {
		for _, printedValue := range indicatorResult.Print() {
			values = append(values, formatValue(printedValue.Value))
		}
	}

	if position != nil {
		values = append(values,
			formatValue(item.StrategyResult.Action),
			formatValue(position.Id),
			formatValue(position.Type),
			formatValue(position.Amount),
			formatValue(position.Sl),
			formatValue(position.Costs),
			formatValue(position.Profit),
			formatValue(position.PossibleProfit),
		)
	}

	return values
}

func writerGetSummaryHeader(summary *Summary) []string {

	var a []string

	for _, value := range summary.StrategyPrintValues {
		a = append(a, value.Label)
	}

	return append(a,
		"Profit",
		"Sum Profitable",
		"Sum Lossy",
		"GrossProfit",
		"Ratio(%)",
		"Positions(+/-)",
		"Avg Positions",
		"Avg Profit",
		"Avg Lost",
	)
}
func writerGetSummaryRow(summary *Summary) []string {
	var a []string

	for _, value := range summary.StrategyPrintValues {
		a = append(a, formatValue(value.Value))
	}

	return append(a,
		formatValue(summary.Profit),
		formatValue(summary.SumOfProfitable),
		formatValue(summary.SumOfLossy),
		formatValue(summary.GrossProfit),
		formatValue(summary.SuccessRatio),
		formatValue(summary.CountOfPositions)+"("+
			formatValue(summary.CountOfProfitable)+"/"+
			formatValue(summary.CountOfLossy)+")",
		formatValue(summary.AvgOfPositions),
		formatValue(summary.AvgOfProfit),
		formatValue(summary.AvgOfLost),
	)
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case float.Float:
		return fmt.Sprintf("%.3f", v.Val())
	case float64:
		return fmt.Sprintf("%.3f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return fmt.Sprintf("%s", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type PrintValue struct {
	Label string
	Value interface{}
}
