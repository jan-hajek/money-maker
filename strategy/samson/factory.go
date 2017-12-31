package samson

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/indicator/adxAvg"
	"github.com/jelito/money-maker/indicator/adxEma"
	"github.com/jelito/money-maker/indicator/sar"
)

type Factory struct {
}

func (s *Factory) GetName() string {
	return "samson"
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) app.StrategyFactoryConfig {
	return Config{
		SarMinimalAf: float.New(config["sarMinimalAf"]["default"].(float64)),
		SarMaximalAf: float.New(config["sarMaximalAf"]["default"].(float64)),
		AdxPeriod:    config["adxPeriod"]["default"].(int),
		SmoothType:   SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha:  float.New(config["smoothAlpha"]["default"].(float64)),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []app.StrategyFactoryConfig {
	sarMinValues := app.FloatSteps(
		float.New(config["sarMinimalAf"]["minimal"].(float64)),
		float.New(config["sarMinimalAf"]["maximal"].(float64)),
		float.New(config["sarMinimalAf"]["step"].(float64)),
	)
	sarMaxValues := app.FloatSteps(
		float.New(config["sarMaximalAf"]["minimal"].(float64)),
		float.New(config["sarMaximalAf"]["maximal"].(float64)),
		float.New(config["sarMaximalAf"]["step"].(float64)),
	)
	adxPeriodValues := app.IntSteps(
		config["adxPeriod"]["minimal"].(int),
		config["adxPeriod"]["maximal"].(int),
		config["adxPeriod"]["step"].(int),
	)

	var smoothTypes []SmoothType
	for _, smoothType := range config["smoothType"]["list"].([]interface{}) {
		smoothTypes = append(smoothTypes, SmoothType(smoothType.(string)))
	}

	smoothAlphaValues := app.FloatSteps(
		float.New(config["smoothAlpha"]["minimal"].(float64)),
		float.New(config["smoothAlpha"]["maximal"].(float64)),
		float.New(config["smoothAlpha"]["step"].(float64)),
	)

	return app.Combinations(
		[]int{
			len(sarMinValues),
			len(sarMaxValues),
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
		},
		func(positions []int) app.StrategyFactoryConfig {
			return Config{
				SarMinimalAf: sarMinValues[positions[0]],
				SarMaximalAf: sarMaxValues[positions[1]],
				AdxPeriod:    adxPeriodValues[positions[2]],
				SmoothType:   smoothTypes[positions[3]],
				SmoothAlpha:  smoothAlphaValues[positions[4]],
			}
		},
	)
}

func (s *Factory) Create(config app.StrategyFactoryConfig) app.Strategy {

	service := &Service{
		config: config.(Config),
		sar: &sar.Service{
			Name:      "sar",
			MinimalAf: config.(Config).SarMinimalAf,
			MaximalAf: config.(Config).SarMaximalAf,
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
			Alpha:  config.(Config).SmoothAlpha,
		}
	}

	return service
}
