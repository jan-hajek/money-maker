package title

import (
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/entity"
)

type Factory interface {
	Create(title *entity.Title) Service
}

type Service interface {
	LoadLast() (dateInput.DateInput, error)
	LoadDataFrom() ([]dateInput.DateInput, error)
}
