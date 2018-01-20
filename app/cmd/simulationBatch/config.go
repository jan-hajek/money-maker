package simulationBatch

type config struct {
	Db         string
	TitleId    string                                       `yaml:"titleId"`
	Strategies map[string]map[string]map[string]interface{} `yaml:"strategies"`
	Writer     struct {
		ParseFormat string `yaml:"parseFormat"`
		Outputs     struct {
			Csv struct {
				Enabled bool
				Path    string
			}
			Stdout struct {
				Enabled bool
			}
		}
	}
}
