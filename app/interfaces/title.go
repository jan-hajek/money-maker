package interfaces

import (
	"github.com/jelito/money-maker/app"
)

type TitleFactory interface {
	Create() TitleService
}

type TitleService interface {
	LoadLast() (app.DateInput, error)
	LoadDataFrom() ([]app.DateInput, error)
}
