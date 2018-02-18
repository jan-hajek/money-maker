package strategy

import (
	"github.com/jelito/money-maker/src/interfaces"
)

type Factory interface {
	GetDefaultConfig(config map[string]map[string]interface{}) FactoryConfig
	GetBatchConfigs(config map[string]map[string]interface{}) []FactoryConfig
	Create(config FactoryConfig) interfaces.Strategy
}

type FactoryConfig interface {
}

func CreateFactoryConfigCombinations(input []int, createResult func([]int) FactoryConfig) []FactoryConfig {
	positions := make([]int, len(input))
	var results []FactoryConfig

	for index := 0; index < len(input); {
		if positions[index] == input[index]-1 {
			index++
		} else {
			results = append(results, createResult(positions))
			positions[index]++
			for y := index - 1; y >= 0; y-- {
				positions[y] = 0
			}
			index = 0
		}
	}

	results = append(results, createResult(positions))

	return results
}
