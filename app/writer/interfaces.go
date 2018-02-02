package writer

import "github.com/jelito/money-maker/app"

type Output interface {
	Open() error
	WriteHistory(historyItems []*app.HistoryItem) error
	WriteSummaryHeader(summary *app.Summary) error
	WriteSummaryRow(summary *app.Summary) error
	Close() error
}

type PrintValue interface {
	GetLabel() string
	GetValue() interface{}
}
