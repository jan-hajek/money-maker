package writer

import (
	"fmt"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/summary"
	"sort"
)

func writerGetHistoryHeader(item interfaces.HistoryItem) []string {
	a := []string{
		"date",
		"price",
	}

	for _, indicatorResult := range orderedIndicatorResults(item) {
		for _, param := range indicatorResult.Print() {
			a = append(a, param.GetLabel())
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
		"poss. prof. %",
	)

	return a
}

func writerGetHistoryRow(item interfaces.HistoryItem, dateFormat string) []string {
	position := item.GetPosition()
	values := []string{
		item.GetDateInput().Date.Format(dateFormat),
		formatValue(item.GetDateInput().ClosePrice),
	}

	for _, indicatorResult := range orderedIndicatorResults(item) {
		for _, printedValue := range indicatorResult.Print() {
			values = append(values, formatValue(printedValue.GetValue()))
		}
	}

	if position != nil {
		values = append(values,
			formatValue(item.GetStrategyResult().GetAction()),
			formatValue(position.Id),
			formatValue(position.Type),
			formatValue(position.Amount),
			formatValue(position.Sl),
			formatValue(position.Costs),
			formatValue(position.Profit),
			formatValue(position.PossibleProfit),
			formatValue(position.PossibleProfitPercent.Val()*100),
		)
	}

	return values
}

func writerGetSummaryHeader(summary *summary.Summary) []string {

	var a []string

	for _, value := range summary.StrategyPrintValues {
		a = append(a, value.GetLabel())
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
func writerGetSummaryRow(summary *summary.Summary) []string {
	var a []string

	for _, value := range summary.StrategyPrintValues {
		a = append(a, formatValue(value.GetValue()))
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

func orderedIndicatorResults(item interfaces.HistoryItem) []interfaces.IndicatorResult {

	var keys []string
	items := item.GetIndicatorResults()
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var ordered []interfaces.IndicatorResult
	for _, k := range keys {
		ordered = append(ordered, items[k])
	}

	return ordered
}
