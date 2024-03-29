package cmd

type config struct {
	Db     string
	Writer struct {
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
	Mail struct {
		Enabled bool
		Addr    string
		From    string
		Pass    string
		To      string
		Host    string
	}
	Run struct {
		DownloadMissingPrices bool `yaml:"downloadMissingPrices"`
	}
	Simulation struct {
		Source struct {
			Csv struct {
				Enabled         bool
				FilePath        string `yaml:"filePath"`
				TimeParseFormat string `yaml:"timeParseFormat"`
			}
			Db struct {
				Enabled bool
				TitleId string `yaml:"titleId"`
			}
		}
		Strategies map[string]map[string]map[string]interface{} `yaml:"strategies"`
	}
}
