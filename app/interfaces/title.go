package interfaces

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
)

type TitleFactory interface {
	Create(title *entity.Title) TitleService
}

type TitleService interface {
	LoadLast() (app.DateInput, error)
	LoadDataFrom() ([]app.DateInput, error)
}
