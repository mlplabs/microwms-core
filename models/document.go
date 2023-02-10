package models

import (
	"database/sql"
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"time"
)

const (
	DocumentTypePosting  = 1
	DocumentTypeReceipt  = 2
	DocumentTypeShipment = 3
	DocumentTypeWriteOff = 4
)

// Document - абстракция документа
type Document struct {
	Name string
	Db   *sql.DB
}

// IDocItem интерфейс документа (пока не используется)
type IDocItem interface {
	GetNumber() string
}

// DocItem строка документа (структура)
type DocItem struct {
	Id      int64    `json:"id"`
	Number  string   `json:"number"`
	Date    string   `json:"date"`
	DocType int      `json:"doc_type"`
	Items   []DocRow `json:"items"`
}

func (d *DocItem) GetNumber() string {
	return fmt.Sprintf("%06d.%d", d.Id, d.DocType)
}

// DocRow товарная строка документа
type DocRow struct {
	RowNum   string  `json:"row_num"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// getItems возвращает список документов (без вложений)
func (d *Document) getItems(offset int, limit int) ([]DocItem, int, error) {
	var count int
	sqlCond := ""

	args := make([]any, 0)

	if limit == 0 {
		limit = 10
	}
	args = append(args, limit)
	args = append(args, offset)

	sqlSel := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s %s ORDER BY date ASC", d.Name, sqlCond)

	rows, err := d.Db.Query(sqlSel+" LIMIT $1 OFFSET $2", args...)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	items := make([]DocItem, count, 10)
	for rows.Next() {
		item := new(DocItem)
		dateDoc := time.Time{}
		err = rows.Scan(&item.Id, &item.Number, &dateDoc, &item.DocType)
		item.Number = item.GetNumber()
		item.Date = GetDate(dateDoc)
		items = append(items, *item)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlSel)
	err = d.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return items, count, nil
}

// createItem создает документ
func (d *Document) createItem(docItem *DocItem) (int64, error) {
	tx, err := d.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlInsH := "INSERT INTO receipt_headers (number, date, doc_type) VALUES($1, $2, $3) RETURNING id"
	err = tx.QueryRow(sqlInsH, docItem.Number, docItem.Date, DocumentTypePosting).Scan(&docItem.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	for idx, item := range docItem.Items {
		if item.Product.Id == 0 {
			if item.Product.Manufacturer.Id == 0 {
				sqlInsMnf := "INSERT INTO manufacturers (name) VALUES ($1) RETURNING id"
				err = tx.QueryRow(sqlInsMnf, item.Product.Manufacturer.Name).Scan(&item.Product.Manufacturer.Id)
				if err != nil {
					tx.Rollback()
					return 0, &core.WrapError{Err: err, Code: core.SystemError}
				}
			}
			sqlInsP := "INSERT INTO products (name, manufacturer_id) VALUES($1, $2) RETURNING id"
			err = tx.QueryRow(sqlInsP, item.Product.Name, item.Product.Manufacturer.Id).Scan(&item.Product.Id)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
		RowId := fmt.Sprintf("%d.%d", docItem.Id, idx+1)
		sqlInsI := "INSERT INTO receipt_items (parent_id, row_id, product_id, quantity) VALUES($1, $2, $3, $4)"
		_, err = tx.Exec(sqlInsI, docItem.Id, RowId, item.Product.Id, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return 0, nil
}

// findItemById находит документ по id (без вложений)
func (d *Document) findItemById(id int64) (*DocItem, error) {
	sqlUsr := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s WHERE id = $1", d.Name)
	row := d.Db.QueryRow(sqlUsr, id)
	u := new(DocItem)
	dateDoc := time.Time{}
	err := row.Scan(&u.Id, &u.Number, &dateDoc, u.DocType)
	u.Number = u.GetNumber()
	u.Date = GetDate(dateDoc)

	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u, nil
}

// getItemById возвращает документ по id (полностью)
func (d *Document) getItemById(id int64) (*DocItem, error) {
	di := DocItem{}
	sqlH := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s WHERE id = $1", d.Name)
	row := d.Db.QueryRow(sqlH, id)
	dateDoc := time.Time{}
	err := row.Scan(&di.Id, &di.Number, dateDoc, &di.DocType)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	di.Number = di.GetNumber()
	di.Date = GetDate(dateDoc)

	sqlI := "SELECT row_id, product_id, quantity FROM receipt_items WHERE parent_id = $1"
	rows, err := d.Db.Query(sqlI, id)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		r := DocRow{}
		err = rows.Scan(&r.RowNum, &r.Product.Id, &r.Quantity)
		di.Items = append(di.Items, r)
	}
	return &di, nil
}

func GetDate(date time.Time) string {
	return date.Format("02.01.2006 15:04:05")
}
