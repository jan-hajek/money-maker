package dateInput

import (
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/float"
)

func CreateFromEntity(pr *entity.Price) app.DateInput {
	return app.DateInput{
		Date:       pr.Date,
		OpenPrice:  float.New(pr.OpenPrice),
		HighPrice:  float.New(pr.HighPrice),
		LowPrice:   float.New(pr.LowPrice),
		ClosePrice: float.New(pr.ClosePrice),
	}
}
