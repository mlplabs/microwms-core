package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
)

type UserService struct {
	Storage *Storage
	Reference
}

// User пользователь
type User struct {
	RefItem
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
	item, err := ps.findItemById(ps.Storage.Db, usrId)
	u := new(User)
	u.RefItem = *item
	return u, err
}

// FindUserByName возвращает пользователя по наименованию
func (ps *UserService) FindUserByName(valName string) ([]User, error) {
	items, err := ps.findItemByName(ps.Storage.Db, valName)
	if err != nil {
		return nil, err
	}
	retVal := make([]User, len(items))
	for _, item := range items {
		u := new(User)
		u.RefItem = item
		retVal = append(retVal, *u)
	}
	return retVal, err
}

func (ps *UserService) GetSuggestionUser(text string, limit int) ([]string, error) {
	return ps.getSuggestion(ps.Storage.Db, text, limit)
}

// CreateUser создает нового пользователя
func (ps *UserService) CreateUser(u *User) (int64, error) {
	return ps.createItem(ps.Storage.Db, u)
}

// DeleteUser удаляет пользователя
func (ps *UserService) DeleteUser(u *User) (int64, error) {
	return ps.deleteItem(ps.Storage.Db, u.Id)
}

// UpdateUser обновляет пользователя
func (ps *UserService) UpdateUser(u *User) (int64, error) {
	return ps.updateItem(ps.Storage.Db, u)
}
