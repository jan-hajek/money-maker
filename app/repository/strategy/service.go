package strategy

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jelito/money-maker/app/entity"
	"github.com/jelito/money-maker/app/interfaces"
)

type Service struct {
	Db *sql.DB
}

func (s *Service) GetById(id string) (*entity.Strategy, error) {
	return s.getRow("select * from strategy WHERE id = ? limit 1", id)
}

func (s *Service) GetAllActive() ([]*entity.Strategy, error) {
	return s.getArray("select * from strategy")
}

func (s *Service) Insert(e *entity.Strategy) error {

	stmt, err := s.Db.Prepare("INSERT INTO strategy(id, name, className) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Id, e.Name, e.ClassName)

	return err
}

func (s *Service) getRow(sql string, args ...interface{}) (*entity.Strategy, error) {
	e := &entity.Strategy{}
	s.scanRow(e, s.Db.QueryRow(sql, args...))

	return e, nil
}
func (s *Service) getArray(sql string) ([]*entity.Strategy, error) {
	rows, err := s.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityList []*entity.Strategy

	for rows.Next() {
		e := &entity.Strategy{}
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

func (s *Service) scanRow(entity *entity.Strategy, result interfaces.Scanable) error {
	return result.Scan(
		&entity.Id,
		&entity.Name,
		&entity.ClassName,
	)
}
