package app

func Combinations(input []int, createResult func([]int) StrategyFactoryConfig) []StrategyFactoryConfig {
	positions := make([]int, len(input))
	var results []StrategyFactoryConfig

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
