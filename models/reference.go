package models

import (
	"database/sql"
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

type Reference struct {
	Name string
}

func (r *Reference) createItem(db *sql.DB, refItem IRefItem) (int64, error) {
	var insertId int64
	sqlCreate := "INSERT INTO users (name) VALUES ($1) RETURNING id"
	err := db.QueryRow(sqlCreate, refItem.GetName()).Scan(&insertId)
	refItem.SetId(insertId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return refItem.GetId(), nil
}

func (r *Reference) updateItem(db *sql.DB, refItem IRefItem) (int64, error) {
	sqlUpd := fmt.Sprintf("UPDATE %s SET name=$2 WHERE id=$1", r.Name)
	res, err := db.Exec(sqlUpd, refItem.GetId(), refItem.GetName())
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	// TODO: корректность возвращаемого значения id или affected?
	return refItem.GetId(), nil
}

func (r *Reference) deleteItem(db *sql.DB, id int64) (int64, error) {
	sqlDel := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.Name)
	res, err := db.Exec(sqlDel, id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}

func (r *Reference) findItemById(db *sql.DB, id int64) (*RefItem, error) {
	sqlUsr := fmt.Sprintf("SELECT id, name FROM %s WHERE m.id = $1", r.Name)
	row := db.QueryRow(sqlUsr, id)
	u := new(RefItem)
	err := row.Scan(&u.Id, &u.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u, nil
}

func (r *Reference) findItemByName(db *sql.DB, name string) ([]RefItem, error) {
	retObjList := make([]RefItem, 0)
	sql := fmt.Sprintf("SELECT id, name FROM %s WHERE name = $1", r.Name)
	rows, err := db.Query(sql, name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		obj := RefItem{}
		err := rows.Scan(&obj.Id, &obj.Name)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retObjList = append(retObjList, obj)
	}
	return retObjList, nil
}

func (r *Reference) getSuggestion(db *sql.DB, text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sqlSel := fmt.Sprintf("SELECT name FROM %s WHERE name LIKE $1 LIMIT $2", r.Name)
	rows, err := db.Query(sqlSel, text+"%", limit)
	if err != nil {
		return retVal, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		s := ""
		err := rows.Scan(&s)
		if err != nil {
			return retVal, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retVal = append(retVal, s)
	}
	return retVal, err
}

type IRefItem interface {
	GetId() int64
	SetId(int64)
	GetName() string
}

type RefItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (r *RefItem) GetId() int64 {
	return r.Id
}

func (r *RefItem) GetName() string {
	return r.Name
}

func (r *RefItem) SetId(id int64) {
	r.Id = id
}
