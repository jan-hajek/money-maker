package plus500

import (
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/title"
)

type Factory struct {
}

func (s *Factory) Create(title *entity.Title) title.Service {
	return &Service{
		titleEntity: title,
	}
}
