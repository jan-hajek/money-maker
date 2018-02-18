package entity

type Trade struct {
	Id         string
	StrategyId string
	TitleId    string
	Params     TradeParams
	Active     int
}
