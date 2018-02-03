package writer

import (
	"fmt"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
)

func writerGetHistoryHeader(item *app.HistoryItem) []string {
	a := []string{
		"date",
		"price",
	}

	for _, indicatorResult := range item.OrderedIndicatorResults() {
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
	)

	return a
}

func writerGetHistoryRow(item *app.HistoryItem, dateFormat string) []string {
	position := item.Position
	values := []string{
		item.DateInput.Date.Format(dateFormat),
		formatValue(item.DateInput.ClosePrice),
	}

	for _, indicatorResult := range item.OrderedIndicatorResults() {
		for _, printedValue := range indicatorResult.Print() {
			values = append(values, formatValue(printedValue.GetValue()))
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

func writerGetSummaryHeader(summary *app.Summary) []string {

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
func writerGetSummaryRow(summary *app.Summary) []string {
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
