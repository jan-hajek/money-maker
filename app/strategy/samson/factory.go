package samson

import (
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEma"
	"github.com/jelito/money-maker/app/indicator/sar"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/strategy"
)

type Factory struct {
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) strategy.FactoryConfig {
	return Config{
		SarMinimalAf: float.New(config["sarMinimalAf"]["default"].(float64)),
		SarMaximalAf: float.New(config["sarMaximalAf"]["default"].(float64)),
		AdxPeriod:    config["adxPeriod"]["default"].(int),
		SmoothType:   SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha:  float.New(config["smoothAlpha"]["default"].(float64)),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []strategy.FactoryConfig {
	sarMinValues := math.FloatSteps(
		float.New(config["sarMinimalAf"]["minimal"].(float64)),
		float.New(config["sarMinimalAf"]["maximal"].(float64)),
		float.New(config["sarMinimalAf"]["step"].(float64)),
	)
	sarMaxValues := math.FloatSteps(
		float.New(config["sarMaximalAf"]["minimal"].(float64)),
		float.New(config["sarMaximalAf"]["maximal"].(float64)),
		float.New(config["sarMaximalAf"]["step"].(float64)),
	)
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

	return strategy.CreateFactoryConfigCombinations(
		[]int{
			len(sarMinValues),
			len(sarMaxValues),
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
		},
		func(positions []int) strategy.FactoryConfig {
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

func (s *Factory) Create(config strategy.FactoryConfig) interfaces.Strategy {

	service := &Service{
		config: config.(Config),
		sar: &sar.Service{
			Name:      "sar",
			MinimalAf: config.(Config).SarMinimalAf,
			MaximalAf: config.(Config).SarMaximalAf,
		},
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
