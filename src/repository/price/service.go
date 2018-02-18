package price

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jelito/money-maker/src/entity"
	"github.com/jelito/money-maker/src/repository"
	"time"
)

type Service struct {
	Db *sql.DB
}

func (s *Service) GetAllActive() ([]*entity.Price, error) {
	return s.getArray("select * from price")
}

func (s *Service) GetLastItemsByTitle(titleId string, limit int) ([]*entity.Price, error) {
	return s.getArray(
		"SELECT sub.* FROM ("+
			"SELECT * FROM price WHERE titleId = ? ORDER BY date DESC LIMIT ?"+
			") AS sub ORDER BY sub.date ASC",
		titleId,
		limit,
	)
}

func (s *Service) GetByTitleAndDate(titleId string, date time.Time) (*entity.Price, error) {
	return s.getRow(
		"SELECT * FROM price WHERE titleId = ? AND date = ? LIMIT 1",
		titleId,
		date,
	)
}

func (s *Service) Insert(e *entity.Price) error {

	stmt, err := s.Db.Prepare("INSERT INTO price(id, titleId, date, openPrice, highPrice, lowPrice, closePrice) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Id, e.TitleId, e.Date, e.OpenPrice, e.HighPrice, e.LowPrice, e.ClosePrice)

	return err
}

func (s *Service) getRow(query string, args ...interface{}) (*entity.Price, error) {
	e := &entity.Price{}
	err := s.scanRow(e, s.Db.QueryRow(query, args...))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return e, err
}

func (s *Service) getArray(sql string, args ...interface{}) ([]*entity.Price, error) {
	rows, err := s.Db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityList []*entity.Price

	for rows.Next() {
		e := &entity.Price{}
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

func (s *Service) scanRow(e *entity.Price, result repository.Scanable) error {
	return result.Scan(
		&e.Id,
		&e.TitleId,
		&e.Date,
		&e.OpenPrice,
		&e.HighPrice,
		&e.LowPrice,
		&e.ClosePrice,
	)
}
