package math

import "github.com/jelito/money-maker/app"

func Combinations(input []int, createResult func([]int) app.ResolverFactoryConfig) []app.ResolverFactoryConfig {
	positions := []int{0, 0, 0}
	var results []app.ResolverFactoryConfig

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
