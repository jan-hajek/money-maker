package entity

import (
	"database/sql"
	"encoding/json"
)

type TradeParams struct {
	sql.NullString `json:"-"`
	Data           map[string]map[string]interface{}
}

func (nt *TradeParams) Scan(value interface{}) error {
	err := nt.NullString.Scan(value)
	if nil != err {
		return err
	}
	if nt.Valid {
		return json.Unmarshal([]byte(nt.String), &nt.Data)
	}
	return nil
}
