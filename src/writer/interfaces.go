package writer

import (
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/summary"
)

type Output interface {
	Open() error
	WriteHistory(historyItems []interfaces.HistoryItem) error
	WriteSummaryHeader(summary *summary.Summary) error
	WriteSummaryRow(summary *summary.Summary) error
	Close() error
}
