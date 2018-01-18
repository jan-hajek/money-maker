package position

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/interfaces"
)

type Service struct {
	Db *sql.DB
}

func (s *Service) LastOpenByTrade(tradeId string) (*entity.Position, error) {
	return s.getRow("SELECT * FROM position WHERE tradeId = ? AND closePriceId IS NULL LIMIT 1", tradeId)
}

func (s *Service) getRow(query string, args ...interface{}) (*entity.Position, error) {
	e := &entity.Position{}
	err := s.scanRow(e, s.Db.QueryRow(query, args...))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return e, err
}

func (s *Service) getArray(sql string) ([]*entity.Position, error) {
	rows, err := s.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityList []*entity.Position

	for rows.Next() {
		e := &entity.Position{}
		err := s.scanRow(e, rows)
		if err != nil {
			return nil, err
		}

		entityList = append(entityList, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return entityList, nil
}

func (s *Service) scanRow(entity *entity.Position, result interfaces.Scanable) error {
	return result.Scan(
		&entity.Id,
		&entity.TradeId,
		&entity.Type,
		&entity.OpenPriceId,
		&entity.ClosePriceId,
		&entity.Amount,
		&entity.Sl,
		&entity.Costs,
		&entity.Profit,
	)
}
