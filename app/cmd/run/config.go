package run

type config struct {
	Db                    string
	DownloadMissingPrices bool `yaml:"downloadMissingPrices"`
	Mail                  struct {
		Enabled bool
		Addr    string
		From    string
		Pass    string
		To      string
	}
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
}
