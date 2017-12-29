package jones

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/float"
	"github.com/jelito/money-maker/indicator/adxAvg"
	"github.com/jelito/money-maker/indicator/adxEma"
)

type Factory struct {
}

func (s *Factory) GetName() string {
	return "jones"
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) app.StrategyFactoryConfig {
	return Config{
		AdxPeriod:   config["adxPeriod"]["default"].(int),
		SmoothType:  SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha: float.New(config["smoothAlpha"]["default"].(float64)),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []app.StrategyFactoryConfig {
	adxPeriodValues := app.IntSteps(
		config["adxPeriod"]["minimal"].(int),
		config["adxPeriod"]["maximal"].(int),
		config["adxPeriod"]["step"].(int),
	)
	smoothTypes := []SmoothType{EMA, AVG}
	smoothAlphaValues := app.FloatSteps(
		float.New(config["smoothAlpha"]["minimal"].(float64)),
		float.New(config["smoothAlpha"]["maximal"].(float64)),
		float.New(config["smoothAlpha"]["step"].(float64)),
	)

	return app.Combinations(
		[]int{
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
		},
		func(positions []int) app.StrategyFactoryConfig {
			return Config{
				AdxPeriod:   adxPeriodValues[positions[0]],
				SmoothType:  smoothTypes[positions[2]],
				SmoothAlpha: smoothAlphaValues[positions[1]],
			}
		},
	)
}

func (s *Factory) Create(config app.StrategyFactoryConfig) app.Strategy {

	service := &Service{
		config: config.(Config),
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
