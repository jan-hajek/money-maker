package title

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jelito/money-maker/app/entity"
)

type Service struct {
	Db *sql.DB
}

func (s *Service) GetAllActive() ([]*entity.Title, error) {
	return s.getArray("select * from title")
}

func (s *Service) getArray(sql string) ([]*entity.Title, error) {
	rows, err := s.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityList []*entity.Title

	for rows.Next() {
		e := entity.Title{}
		err := rows.Scan(&e.Id, &e.Name, &e.DataUrl, &e.DownloadInterval, &e.ClassName)
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
