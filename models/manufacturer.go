package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

// Manufacturer производитель
type Manufacturer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// CreateManufacturer создает новый продукт
func (ps *ProductService) CreateManufacturer(m *Manufacturer) (int64, error) {
	sqlInsProd := "INSERT INTO manufacturers (name) VALUES ($1) RETURNING id"
	err := ps.Storage.Db.QueryRow(sqlInsProd, m.Name).Scan(&m.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return m.Id, nil
}

// GetManufacturers возвращает список производителей
func (ps *ProductService) GetManufacturers(offset int, limit int) ([]Manufacturer, int, error) {
	var count int

	sqlMnf := "SELECT m.id, m.name FROM manufacturers m " +
		"		ORDER BY m.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := ps.Storage.Query(sqlMnf+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	mnfs := make([]Manufacturer, count, 10)
	for rows.Next() {
		m := new(Manufacturer)
		err = rows.Scan(&m.Id, &m.Name)
		mnfs = append(mnfs, *m)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlMnf)
	err = ps.Storage.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return mnfs, count, nil
}

// FindManufacturerById возвращает производителя по внутреннему идентификатору
func (ps *ProductService) FindManufacturerById(mnfId int64) (*Manufacturer, error) {

	sqlCell := "SELECT m.id, m.name " +
		"FROM manufacturers m " +
		"WHERE m.id = $1"
	row := ps.Storage.Db.QueryRow(sqlCell, mnfId)
	m := new(Manufacturer)
	err := row.Scan(&m.Id, &m.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return m, nil
}

// FindManufacturerByName возвращает список производителей по наименованию
func (ps *ProductService) FindManufacturerByName(mnfName string) ([]Manufacturer, error) {
	retMnf := make([]Manufacturer, 0)
	sql := "SELECT id, name FROM manufacturers WHERE name = $1"
	rows, err := ps.Storage.Query(sql, mnfName)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		m := Manufacturer{}
		err := rows.Scan(&m.Id, &m.Name)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retMnf = append(retMnf, m)
	}
	return retMnf, nil
}

// UpdateManufacturer создает новый продукт
func (ps *ProductService) UpdateManufacturer(m *Manufacturer) (int64, error) {
	sqlUpd := "UPDATE manufacturers SET name=$2 WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlUpd, m.Id, m.Name)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return m.Id, nil
}

func (ps *ProductService) GetSuggestionManufacturers(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sql := "SELECT name FROM manufacturers WHERE name LIKE $1 LIMIT $2"
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

func (ps *ProductService) DeleteManufacturer(m *Manufacturer) (int64, error) {
	sqlDel := "DELETE FROM manufacturers WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlDel, m.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
