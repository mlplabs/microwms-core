package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

type UserService struct {
	Storage *Storage
}

// User пользователь
type User struct {
	Reference
}

// CreateUser создает нового пользователя
func (ps *UserService) CreateUser(u *User) (int64, error) {
	return u.createItem(ps.Storage.Db)
}

// GetUsers возвращает список пользователей
func (ps *UserService) GetUsers(offset int, limit int) ([]User, int, error) {
	var count int

	sqlUsr := "SELECT m.id, m.name FROM users m " +
		"		ORDER BY m.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := ps.Storage.Query(sqlUsr+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	usrs := make([]User, count, 10)
	for rows.Next() {
		u := new(User)
		err = rows.Scan(&u.Id, &u.Name)
		usrs = append(usrs, *u)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlUsr)
	err = ps.Storage.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return usrs, count, nil
}

// FindUserById возвращает пользователя по внутреннему идентификатору
func (ps *UserService) FindUserById(usrId int64) (*User, error) {

	sqlUsr := "SELECT m.id, m.name " +
		"FROM users m " +
		"WHERE m.id = $1"
	row := ps.Storage.Db.QueryRow(sqlUsr, usrId)
	u := new(User)
	err := row.Scan(&u.Id, &u.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u, nil
}

// FindUserByName возвращает пользователя по наименованию
func (ps *UserService) FindUserByName(valName string) ([]User, error) {
	retObjList := make([]User, 0)
	sql := "SELECT id, name FROM users WHERE name = $1"
	rows, err := ps.Storage.Query(sql, valName)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		obj := User{}
		err := rows.Scan(&obj.Id, &obj.Name)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retObjList = append(retObjList, obj)
	}
	return retObjList, nil
}

// UpdateUser обновляет пользователя
func (ps *UserService) UpdateUser(u *User) (int64, error) {
	sqlUpd := "UPDATE users SET name=$2 WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlUpd, u.Id, u.Name)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u.Id, nil
}

func (ps *UserService) GetSuggestionUser(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sql := "SELECT name FROM users WHERE name LIKE $1 LIMIT $2"
	rows, err := ps.Storage.Query(sql, text+"%", limit)
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

func (ps *UserService) DeleteUser(u *User) (int64, error) {
	sqlDel := "DELETE FROM manufacturers WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlDel, u.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
