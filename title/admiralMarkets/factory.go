package admiralMarkets

import "github.com/jelito/money-maker/app/interfaces"

type Factory struct {
}

func (s *Factory) Create() interfaces.TitleService {
	return &Service{}
}
