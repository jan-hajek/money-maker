package writer

import (
	"fmt"
	"github.com/jelito/money-maker/app"
	"io"
	"os"
	"text/tabwriter"
)

type StdOutWriterOutput struct {
	DateFormat string
	w          *tabwriter.Writer
}

func (s *StdOutWriterOutput) Open() error {
	s.w = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)

	return nil
}

func (s *StdOutWriterOutput) WriteHistory(historyItems []*app.HistoryItem) error {
	if len(historyItems) == 0 {
		return nil
	}
	s.write(s.w, writerGetHistoryHeader(historyItems[0])...)
	for _, item := range historyItems {
		s.write(s.w, writerGetHistoryRow(item, s.DateFormat)...)
	}
	return nil
}

func (s *StdOutWriterOutput) WriteSummaryHeader(summary *app.Summary) error {
	s.write(s.w, writerGetSummaryHeader(summary)...)

	return nil
}

func (s *StdOutWriterOutput) WriteSummaryRow(summary *app.Summary) error {
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
