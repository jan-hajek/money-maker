package jones2

import (
	"github.com/jelito/money-maker/app/indicator/adxAvg"
	"github.com/jelito/money-maker/app/indicator/adxEmaRSI"
	"github.com/jelito/money-maker/src/interfaces"
	"github.com/jelito/money-maker/src/math"
	"github.com/jelito/money-maker/src/math/float"
	"github.com/jelito/money-maker/src/strategy"
)

type Factory struct {
}

func (s *Factory) GetName() string {
	return "jones2"
}

func (s *Factory) GetDefaultConfig(config map[string]map[string]interface{}) strategy.FactoryConfig {
	return Config{
		AdxPeriod:     int(config["adxPeriod"]["default"].(float64)),
		SmoothType:    SmoothType(config["smoothType"]["default"].(string)),
		SmoothAlpha:   float.New(config["smoothAlpha"]["default"].(float64)),
		OpenLowerAdx:  int(config["openLowerAdx"]["default"].(float64)),
		OpenHigherAdx: int(config["openHigherAdx"]["default"].(float64)),
		CloseAdx:      int(config["closeAdx"]["default"].(float64)),
		DIOpenLevel:   int(config["diOpenLevel"]["default"].(float64)),
		DICloseLevel:  int(config["diCloseLevel"]["default"].(float64)),
		DISDCount:     float.New(config["diSDcount"]["default"].(float64)),
		PeriodDIMA:    int(config["periodDIMA"]["default"].(float64)),
		Spread:        float.New(config["spread"]["default"].(float64)),
		Swap:          float.New(config["swap"]["default"].(float64)),
		StopProfit:    float.New(config["stopProfit"]["default"].(float64)),
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
		config["openLowerAdx"]["minimal"].(int),
		config["openLowerAdx"]["maximal"].(int),
		config["openLowerAdx"]["step"].(int),
	)

	openHigherAdxValues := math.IntSteps(
		config["openHigherAdx"]["minimal"].(int),
		config["openHigherAdx"]["maximal"].(int),
		config["openHigherAdx"]["step"].(int),
	)

	closeAdxValues := math.IntSteps(
		config["closeAdx"]["minimal"].(int),
		config["closeAdx"]["maximal"].(int),
		config["closeAdx"]["step"].(int),
	)

	diOpenLevelValues := math.IntSteps(
		config["diOpenLevel"]["minimal"].(int),
		config["diOpenLevel"]["maximal"].(int),
		config["diOpenLevel"]["step"].(int),
	)

	diCloseLevelValues := math.IntSteps(
		config["diCloseLevel"]["minimal"].(int),
		config["diCloseLevel"]["maximal"].(int),
		config["diCloseLevel"]["step"].(int),
	)

	diSDCountValues := math.FloatSteps(
		float.New(config["diSDcount"]["minimal"].(float64)),
		float.New(config["diSDcount"]["maximal"].(float64)),
		float.New(config["diSDcount"]["step"].(float64)),
	)
	periodDIMAValues := math.IntSteps(
		config["periodDIMA"]["minimal"].(int),
		config["periodDIMA"]["maximal"].(int),
		config["periodDIMA"]["step"].(int),
	)
	spreadValues := math.FloatSteps(
		float.New(config["spread"]["minimal"].(float64)),
		float.New(config["spread"]["maximal"].(float64)),
		float.New(config["spread"]["step"].(float64)),
	)
	swapValues := math.FloatSteps(
		float.New(config["swap"]["minimal"].(float64)),
		float.New(config["swap"]["maximal"].(float64)),
		float.New(config["swap"]["step"].(float64)),
	)
	stopProfitValues := math.FloatSteps(
		float.New(config["stopProfit"]["minimal"].(float64)),
		float.New(config["stopProfit"]["maximal"].(float64)),
		float.New(config["stopProfit"]["step"].(float64)),
	)

	return strategy.CreateFactoryConfigCombinations(
		[]int{
			len(adxPeriodValues),
			len(smoothTypes),
			len(smoothAlphaValues),
			len(openLowerAdxValues),
			len(openHigherAdxValues),
			len(closeAdxValues),
			len(diOpenLevelValues),
			len(diCloseLevelValues),
			len(diSDCountValues),
			len(periodDIMAValues),
			len(spreadValues),
			len(swapValues),
			len(stopProfitValues),
		},
		func(positions []int) strategy.FactoryConfig {
			return Config{
				AdxPeriod:     adxPeriodValues[positions[0]],
				SmoothType:    smoothTypes[positions[1]],
				SmoothAlpha:   smoothAlphaValues[positions[2]],
				OpenLowerAdx:  openLowerAdxValues[positions[3]],
				OpenHigherAdx: openHigherAdxValues[positions[4]],
				CloseAdx:      closeAdxValues[positions[5]],
				DIOpenLevel:   diOpenLevelValues[positions[6]],
				DICloseLevel:  diCloseLevelValues[positions[7]],
				DISDCount:     diSDCountValues[positions[8]],
				PeriodDIMA:    periodDIMAValues[positions[9]],
				Spread:        spreadValues[positions[10]],
				Swap:          swapValues[positions[11]],
				StopProfit:    stopProfitValues[positions[12]],
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
		service.adxEmaRSI = &adxEmaRSI.Service{
			Name:       "adx",
			Period:     config.(Config).AdxPeriod,
			Alpha:      config.(Config).SmoothAlpha,
			PeriodDIMA: config.(Config).PeriodDIMA,
			DISDCount:  config.(Config).DISDCount,
		}
	}

	return service
}
