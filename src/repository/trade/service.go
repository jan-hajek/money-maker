package trade

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jelito/money-maker/src/entity"
)

type Service struct {
	Db *sql.DB
}

func (s *Service) GetAllActive() ([]*entity.Trade, error) {
	return s.getArray("select * from trade where active = 1")
}

func (s *Service) GetAllActiveByTitleId(titleId string) ([]*entity.Trade, error) {
	return s.getArray("select * from trade where titleId = ?", titleId)
}

func (s *Service) getArray(sql string, args ...interface{}) ([]*entity.Trade, error) {
	rows, err := s.Db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityList []*entity.Trade

	for rows.Next() {
		e := entity.Trade{}
		err := rows.Scan(
			&e.Id,
			&e.StrategyId,
			&e.TitleId,
			&e.Params,
			&e.Active,
		)
		if err != nil {
			return nil, err
		}

		entityList = append(entityList, &e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return entityList, nil
}
