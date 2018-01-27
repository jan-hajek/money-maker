package dateInput

import (
	"github.com/jelito/money-maker/app"
)

type Loader interface {
	Load() ([]app.DateInput, error)
}
