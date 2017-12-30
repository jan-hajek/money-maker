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
		AdxPeriod:        config["adxPeriod"]["default"].(int),
		SmoothType:       SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha:      float.New(config["smoothAlpha"]["default"].(float64)),
		OpenLowerAdxVal:  config["openLowerAdxVal"]["default"].(int),
		OpenHigherAdxVal: config["openHigherAdxVal"]["default"].(int),
		CloseAdxVal:      config["closeAdxVal"]["default"].(int),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []app.StrategyFactoryConfig {
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

	openLowerAdxValues := app.IntSteps(
		config["openLowerAdxVal"]["minimal"].(int),
		config["openLowerAdxVal"]["maximal"].(int),
		config["openLowerAdxVal"]["step"].(int),
	)

	openHigherAdxValues := app.IntSteps(
		config["openHigherAdxVal"]["minimal"].(int),
		config["openHigherAdxVal"]["maximal"].(int),
		config["openHigherAdxVal"]["step"].(int),
	)

	closeAdxValues := app.IntSteps(
		config["closeAdxVal"]["minimal"].(int),
		config["closeAdxVal"]["maximal"].(int),
		config["closeAdxVal"]["step"].(int),
	)

	return app.Combinations(
		[]int{
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
			len(openLowerAdxValues),
			len(openHigherAdxValues),
			len(closeAdxValues),
		},
		func(positions []int) app.StrategyFactoryConfig {
			return Config{
				AdxPeriod:        adxPeriodValues[positions[0]],
				SmoothType:       smoothTypes[positions[1]],
				SmoothAlpha:      smoothAlphaValues[positions[2]],
				OpenLowerAdxVal:  openLowerAdxValues[positions[3]],
				OpenHigherAdxVal: openHigherAdxValues[positions[4]],
				CloseAdxVal:      closeAdxValues[positions[5]],
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
			Alpha:  config.(Config).SmoothAlpha,
		}
	}

	return service
}
