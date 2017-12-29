package samson

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/calculator/adxAvg"
	"github.com/jelito/money-maker/calculator/adxEma"
	"github.com/jelito/money-maker/calculator/sar"
)

type Factory struct {
}

func (s *Factory) GetName() string {
	return "samson"
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) app.StrategyFactoryConfig {
	return Config{
		SarMinimalAf: config["sarMinimalAf"]["default"].(float64),
		SarMaximalAf: config["sarMaximalAf"]["default"].(float64),
		AdxPeriod:    config["adxPeriod"]["default"].(int),
		SmoothType:   SmoothType(config["smoothType"]["default"].(string)),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []app.StrategyFactoryConfig {
	sarMinValues := app.FloatSteps(
		config["sarMinimalAf"]["minimal"].(float64),
		config["sarMinimalAf"]["maximal"].(float64),
		config["sarMinimalAf"]["step"].(float64),
	)
	sarMaxValues := app.FloatSteps(
		config["sarMaximalAf"]["minimal"].(float64),
		config["sarMaximalAf"]["maximal"].(float64),
		config["sarMaximalAf"]["step"].(float64),
	)
	adxPeriodValues := app.IntSteps(
		config["adxPeriod"]["minimal"].(int),
		config["adxPeriod"]["maximal"].(int),
		config["adxPeriod"]["step"].(int),
	)
	smoothTypes := []SmoothType{EMA, AVG}

	return app.Combinations(
		[]int{
			len(sarMinValues),
			len(sarMaxValues),
			len(adxPeriodValues),
			len(smoothTypes),
		},
		func(positions []int) app.StrategyFactoryConfig {
			return Config{
				SarMinimalAf: sarMinValues[positions[0]],
				SarMaximalAf: sarMaxValues[positions[1]],
				AdxPeriod:    adxPeriodValues[positions[2]],
				SmoothType:   smoothTypes[positions[3]],
			}
		},
	)
}

func (s *Factory) Create(config app.StrategyFactoryConfig) app.Strategy {

	service := &Service{
		config: config.(Config),
		sar: &sar.Service{
			Name:      "sar",
			MinimalAf: float.New(config.(Config).SarMinimalAf),
			MaximalAf: float.New(config.(Config).SarMaximalAf),
		},
	}

	if service.config.SmoothType == AVG {
		service.adxAvg = &adxAvg.Service{
			Name:   "adx",
			Period: config.(Config).AdxPeriod,
		}
	}

	if service.config.SmoothType == EMA {
		service.adxEma = &adxEma.Service{
			Name:   "adx",
			Period: config.(Config).AdxPeriod,
		}
	}

	return service
}
