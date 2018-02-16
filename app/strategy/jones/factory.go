package jones

import (
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEma"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/strategy"
)

type Factory struct {
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) strategy.FactoryConfig {
	return Config{
		AdxPeriod:        int(config["adxPeriod"]["default"].(float64)),
		SmoothType:       SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha:      float.New(config["smoothAlpha"]["default"].(float64)),
		OpenLowerAdxVal:  int(config["openLowerAdxVal"]["default"].(float64)),
		OpenHigherAdxVal: int(config["openHigherAdxVal"]["default"].(float64)),
		CloseAdxVal:      int(config["closeAdxVal"]["default"].(float64)),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []strategy.FactoryConfig {
	adxPeriodValues := math.IntSteps(
		config["adxPeriod"]["minimal"].(int),
		config["adxPeriod"]["maximal"].(int),
		config["adxPeriod"]["step"].(int),
	)

	var smoothTypes []SmoothType
	for _, smoothType := range config["smoothType"]["list"].([]interface{}) {
		smoothTypes = append(smoothTypes, SmoothType(smoothType.(string)))
	}

	smoothAlphaValues := math.FloatSteps(
		float.New(config["smoothAlpha"]["minimal"].(float64)),
		float.New(config["smoothAlpha"]["maximal"].(float64)),
		float.New(config["smoothAlpha"]["step"].(float64)),
	)

	openLowerAdxValues := math.IntSteps(
		config["openLowerAdxVal"]["minimal"].(int),
		config["openLowerAdxVal"]["maximal"].(int),
		config["openLowerAdxVal"]["step"].(int),
	)

	openHigherAdxValues := math.IntSteps(
		config["openHigherAdxVal"]["minimal"].(int),
		config["openHigherAdxVal"]["maximal"].(int),
		config["openHigherAdxVal"]["step"].(int),
	)

	closeAdxValues := math.IntSteps(
		config["closeAdxVal"]["minimal"].(int),
		config["closeAdxVal"]["maximal"].(int),
		config["closeAdxVal"]["step"].(int),
	)

	return strategy.CreateFactoryConfigCombinations(
		[]int{
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
			len(openLowerAdxValues),
			len(openHigherAdxValues),
			len(closeAdxValues),
		},
		func(positions []int) strategy.FactoryConfig {
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

func (s *Factory) Create(config strategy.FactoryConfig) interfaces.Strategy {

	service := &Service{
		config: config.(Config),
	}

	if service.config.SmoothType == AVG {
		service.adxAvg = adxAvg.New(
			"adx",
			config.(Config).AdxPeriod,
		)
	}

	if service.config.SmoothType == EMA {
		service.adxEma = adxEma.New(
			"adx",
			config.(Config).AdxPeriod,
			config.(Config).SmoothAlpha,
		)
	}

	return service
}
