package models

import (
	"database/sql"
	"github.com/mlplabs/microwms-core/core"
)

type Reference struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (r *Reference) createItem(db *sql.DB) (int64, error) {
	sqlCreate := "INSERT INTO users (name) VALUES ($1) RETURNING id"
	err := db.QueryRow(sqlCreate, r.Name).Scan(&r.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return r.Id, nil
}

func (r *Reference) deleteItem(db *sql.DB) (int64, error) {
	sqlDel := "DELETE FROM users WHERE id=$1"
	res, err := db.Exec(sqlDel, r.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
