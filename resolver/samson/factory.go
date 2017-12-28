package samson

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/math"
	"github.com/jelito/money-maker/calculator"
)

type Factory struct {
}

func (s *Factory) GetName() string {
	return "samson"
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) app.ResolverFactoryConfig {
	return Config{
		SarMinimalAf: config["sarMinimalAf"]["default"].(float64),
		SarMaximalAf: config["sarMaximalAf"]["default"].(float64),
		AdxPeriod:    config["adxPeriod"]["default"].(int),
	}
}

func (s *Factory) GetBatchConfigs(config map[string]map[string]interface{}) []app.ResolverFactoryConfig {
	sarMinValues := math.FloatSteps(
		config["sarMinimalAf"]["minimal"].(float64),
		config["sarMinimalAf"]["maximal"].(float64),
		config["sarMinimalAf"]["step"].(float64),
	)
	sarMaxValues := math.FloatSteps(
		config["sarMaximalAf"]["minimal"].(float64),
		config["sarMaximalAf"]["maximal"].(float64),
		config["sarMaximalAf"]["step"].(float64),
	)
	adxPeriodValues := math.IntSteps(
		config["adxPeriod"]["minimal"].(int),
		config["adxPeriod"]["maximal"].(int),
		config["adxPeriod"]["step"].(int),
	)

	ccc := math.Combinations(
		[]int{
			len(sarMinValues),
			len(sarMaxValues),
			len(adxPeriodValues),
		},
		func(positions []int) app.ResolverFactoryConfig {
			return Config{
				SarMinimalAf: sarMinValues[positions[0]],
				SarMaximalAf: sarMaxValues[positions[1]],
				AdxPeriod:    adxPeriodValues[positions[2]],
			}
		},
	)

	return ccc
}

func (s *Factory) Create(config app.ResolverFactoryConfig) app.Resolver {

	return &Service{
		config: config.(Config),
		sar: &calculator.Sar{
			Name:      "sar",
			MinimalAf: config.(Config).SarMinimalAf,
			MaximalAf: config.(Config).SarMaximalAf,
		},

		adx: &calculator.Adx{
			Name:   "adx",
			Period: config.(Config).AdxPeriod,
		},
	}
}
