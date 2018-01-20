package cmd

import (
	"database/sql"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/position"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/repository/strategy"
	"github.com/jelito/money-maker/app/repository/title"
	"github.com/jelito/money-maker/app/repository/trade"
	"github.com/jelito/money-maker/strategy/jones"
	"github.com/jelito/money-maker/strategy/jones2"
	"github.com/jelito/money-maker/strategy/samson"
	"github.com/jelito/money-maker/title/admiralMarkets"
)

func AddDefaultClasses(reg *registry.Registry) {
	db := reg.GetByName("db").(*sql.DB)

	reg.Add("app/repository/trade", &trade.Service{Db: db})
	reg.Add("app/repository/title", &title.Service{Db: db})
	reg.Add("app/repository/price", &price.Service{Db: db})
	reg.Add("app/repository/strategy", &strategy.Service{Db: db})
	reg.Add("app/repository/position", &position.Service{Db: db})

	reg.Add("strategy/samson", &samson.Factory{})
	reg.Add("strategy/jones", &jones.Factory{})
	reg.Add("strategy/jones2", &jones2.Factory{})
	reg.Add("title/admiralMarkets", &admiralMarkets.Factory{})
}
