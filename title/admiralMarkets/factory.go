package admiralMarkets

import (
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/interfaces"
)

type Factory struct {
}

func (s *Factory) Create(title *entity.Title) interfaces.TitleService {
	return &Service{
		titleEntity: title,
	}
}
