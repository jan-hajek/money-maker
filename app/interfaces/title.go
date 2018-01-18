package interfaces

import (
	"github.com/jelito/money-maker/app"
)

type TitleFactory interface {
	Create() TitleService
}

type TitleService interface {
	LoadLast() app.DateInput
	LoadDataFrom() []app.DateInput
}
