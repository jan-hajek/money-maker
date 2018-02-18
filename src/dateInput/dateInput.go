package dateInput

import (
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/math/float"
	"time"
)

type DateInput struct {
	Date       time.Time
	OpenPrice  float.Float
	ClosePrice float.Float
	HighPrice  float.Float
	LowPrice   float.Float
}

func CreateFromEntity(pr *entity.Price) DateInput {
	return DateInput{
		Date:       pr.Date,
		OpenPrice:  float.New(pr.OpenPrice),
		HighPrice:  float.New(pr.HighPrice),
		LowPrice:   float.New(pr.LowPrice),
		ClosePrice: float.New(pr.ClosePrice),
	}
}
