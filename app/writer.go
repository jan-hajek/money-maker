package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type Writer struct {
	outputs []WriterOutput
}

func (s *Writer) Open() error {
	for _, output := range s.outputs {
		err := output.Open()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteHistory(history *History) error {
	for _, output := range s.outputs {
		err := output.WriteHistory(history)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) WriteSummary(summaryList []*Summary) error {
	for _, output := range s.outputs {
		err := output.WriteSummary(summaryList)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Writer) Close() error {
	for _, output := range s.outputs {
		err := output.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

type WriterOutput interface {
	Open() error
	WriteHistory(history *History) error
	WriteSummary(summaryList []*Summary) error
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

func (s *StdOutWriterOutput) WriteHistory(history *History) error {
	s.write(s.w, writerGetHistoryHeader(history)...)

	for _, item := range history.GetAll() {
		s.write(s.w, writerGetHistoryRow(item, s.DateFormat)...)
	}

	return nil
}

func (s *StdOutWriterOutput) WriteSummary(summaryList []*Summary) error {
	s.write(s.w, writerGetSummaryHeader(summaryList[0])...)
	for _, summary := range summaryList {
		s.write(s.w, writerGetSummaryRow(summary)...)
	}
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
	s.file, err = os.Create(s.File)
	if err != nil {
		return err
	}

	s.writer = csv.NewWriter(s.file)
	s.writer.Comma = ';'

	return err
}

func (s *CsvWriterOutput) WriteHistory(history *History) error {

	s.writer.Write(writerGetHistoryHeader(history))

	for _, item := range history.GetAll() {
		s.writer.Write(writerGetHistoryRow(item, s.DateFormat))
	}

	return nil
}

func (s *CsvWriterOutput) WriteSummary(summaryList []*Summary) error {
	s.writer.Write(writerGetSummaryHeader(summaryList[0]))
	for _, summary := range summaryList {
		s.writer.Write(writerGetSummaryRow(summary))
	}
	return nil
}

func (s *CsvWriterOutput) Close() error {
	s.writer.Flush()
	return s.file.Close()
}

func writerGetHistoryHeader(history *History) []string {
	a := []string{
		"date",
		"price",
		"",
	}

	for _, item := range history.GetLastItems(1) {
		for _, calculatorResult := range item.OrderedCalculatorResults() {
			for _, param := range calculatorResult.Print() {
				a = append(a, param.Label)
			}
		}
	}

	a = append(a,
		"",
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
		fmt.Sprintf("%.3f", item.DateInput.ClosePrice),
		"",
	}

	for _, calculatorResult := range item.OrderedCalculatorResults() {
		for _, printedValue := range calculatorResult.Print() {
			values = append(values, formatValue(printedValue.Value))
		}
	}

	if position != nil {
		values = append(values,
			"",
			fmt.Sprintf("%s", item.ResolverResult.Action),
			fmt.Sprintf("%d", position.Id),
			fmt.Sprintf("%s", position.Type),
			fmt.Sprintf("%.3f", position.Amount),
			fmt.Sprintf("%.3f", position.Sl),
			fmt.Sprintf("%.3f", position.Costs),
			fmt.Sprintf("%.3f", position.Profit),
			fmt.Sprintf("%.3f", position.PossibleProfit),
		)
	}

	return values
}

func writerGetSummaryHeader(summary *Summary) []string {

	var a []string

	for _, value := range summary.ResolverPrintValues {
		a = append(a, value.Label)
	}

	return append(a,
		"Profit",
		"Sum Profitable",
		"Sum Lossy",
		"Positions",
		"Profitable",
		"Lossy",
		"Avg Positions",
		"Avg Profitable",
		"Avg Lossy",
	)
}
func writerGetSummaryRow(summary *Summary) []string {
	var a []string

	for _, value := range summary.ResolverPrintValues {
		a = append(a, formatValue(value.Value))
	}

	return append(a,
		fmt.Sprintf("%.3f", summary.Profit),
		fmt.Sprintf("%.3f", summary.SumOfProfitable),
		fmt.Sprintf("%.3f", summary.SumOfLossy),
		fmt.Sprintf("%d", summary.CountOfPositions),
		fmt.Sprintf("%d", summary.CountOfProfitable),
		fmt.Sprintf("%d", summary.CountOfLossy),
		fmt.Sprintf("%.3f", summary.AvgOfPositions),
		fmt.Sprintf("%.3f", summary.AvgOfProfitable),
		fmt.Sprintf("%.3f", summary.AvgOfLossy),
	)
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
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
