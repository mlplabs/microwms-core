package whs

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
	HeadersName string
	ItemsName   string
	Db          *sql.DB
}

// IDocItem интерфейс документа
type IDocItem interface {
	GetNumber() string
	getId() int64
	getType() int
}

type Doc struct {
	Id      int64  `json:"id"`
	Number  string `json:"number"`
	Date    string `json:"date"`
	DocType int    `json:"doc_type"`
}

// DocItem строка документа (структура)
type DocItem struct {
	Doc
	Items []DocRow `json:"items"`
}

func (d *DocItem) getId() int64 {
	return d.Id
}

func (d *DocItem) getType() int {
	return d.DocType
}

func (d *DocItem) GetNumber() string {
	return fmt.Sprintf("%06d.%d", d.Id, d.DocType)
}

// DocRow product line of the document
type DocRow struct {
	RowId    string  `json:"row_id"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	CellSrc  Cell    `json:"cell_src"` // from
	CellDst  Cell    `json:"cell_dst"` // to
}

// getItems returns a list of documents (without goods)
func (d *Document) getItems(offset int, limit int) ([]DocItem, int, error) {
	var count int
	sqlCond := ""

	args := make([]any, 0)

	if limit == 0 {
		limit = 10
	}
	args = append(args, limit)
	args = append(args, offset)

	sqlSel := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s %s ORDER BY date ASC", d.HeadersName, sqlCond)

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

// createItem creates a document
func (d *Document) createItem(docItem *DocItem) (int64, error) {
	tx, err := d.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlInsH := fmt.Sprintf("INSERT INTO %s (number, date, doc_type) VALUES($1, $2, $3) RETURNING id", d.HeadersName)
	err = tx.QueryRow(sqlInsH, docItem.Number, docItem.Date, DocumentTypeReceipt).Scan(&docItem.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	for idx, item := range docItem.Items {
		if item.Product.Id == 0 {
			if item.Product.Manufacturer.Id == 0 {
				sqlInsMnf := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", tableRefManufacturers)
				err = tx.QueryRow(sqlInsMnf, item.Product.Manufacturer.Name).Scan(&item.Product.Manufacturer.Id)
				if err != nil {
					tx.Rollback()
					return 0, &core.WrapError{Err: err, Code: core.SystemError}
				}
			}
			sqlInsP := fmt.Sprintf("INSERT INTO %s (name, manufacturer_id) VALUES($1, $2) RETURNING id", tableRefProducts)
			err = tx.QueryRow(sqlInsP, item.Product.Name, item.Product.Manufacturer.Id).Scan(&item.Product.Id)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
		item.RowId = fmt.Sprintf("%d.%d", docItem.Id, idx+1)
		sqlInsI := fmt.Sprintf("INSERT INTO %s (parent_id, row_id, product_id, quantity) VALUES($1, $2, $3, $4)", d.ItemsName)
		_, err = tx.Exec(sqlInsI, docItem.Id, item.RowId, item.Product.Id, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}

		c := Cell{Id: 2, WhsId: 1, ZoneId: 1}
		s := Storage{Db: d.Db}
		item.CellDst = c

		_, err = s.PutRow(docItem, &item, tx)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}

		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return docItem.Id, nil
}

// updateItem updates the document
func (d *Document) updateItem(docItem *DocItem) (int64, error) {
	tx, err := d.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlUpdH := fmt.Sprintf("UPDATE %s SET number = $1, date = $2, doc_type = $3 WHERE id = $4", d.HeadersName)
	_, err = tx.Exec(sqlUpdH, docItem.Number, docItem.Date, DocumentTypePosting, docItem.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlDelProd := fmt.Sprintf("DELETE FROM %s WHERE parent_id = $1", d.ItemsName)
	_, err = tx.Exec(sqlDelProd, docItem.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for idx, item := range docItem.Items {
		if item.Product.Id == 0 {
			if item.Product.Manufacturer.Id == 0 {
				sqlInsMnf := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", tableRefManufacturers)
				err = tx.QueryRow(sqlInsMnf, item.Product.Manufacturer.Name).Scan(&item.Product.Manufacturer.Id)
				if err != nil {
					tx.Rollback()
					return 0, &core.WrapError{Err: err, Code: core.SystemError}
				}
			}
			sqlInsP := fmt.Sprintf("INSERT INTO %s (name, manufacturer_id) VALUES($1, $2) RETURNING id", tableRefProducts)
			err = tx.QueryRow(sqlInsP, item.Product.Name, item.Product.Manufacturer.Id).Scan(&item.Product.Id)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
		RowId := fmt.Sprintf("%d.%d", docItem.Id, idx+1)
		sqlInsI := fmt.Sprintf("INSERT INTO %s (parent_id, row_id, product_id, quantity) VALUES($1, $2, $3, $4)", d.ItemsName)
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

	return docItem.Id, nil
}

// findItemById finds a document by id (without goods)
func (d *Document) findItemById(id int64) (*DocItem, error) {
	sqlUsr := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s WHERE id = $1", d.HeadersName)
	row := d.Db.QueryRow(sqlUsr, id)
	u := new(DocItem)
	dateDoc := time.Time{}
	err := row.Scan(&u.Id, &u.Number, &dateDoc, &u.DocType)
	u.Number = u.GetNumber()
	u.Date = GetDate(dateDoc)

	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u, nil
}

// findItemByNumberDate finds a document by number and date (without goods)
func (d *Document) findItemByNumberDate(number string, date time.Time) (*DocItem, error) {
	sqlUsr := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s WHERE number = $1 AND date::date >= $2::date", d.HeadersName)
	row := d.Db.QueryRow(sqlUsr, number, date)
	u := new(DocItem)
	dateDoc := time.Time{}
	err := row.Scan(&u.Id, &u.Number, &dateDoc, &u.DocType)
	u.Number = u.GetNumber()
	u.Date = GetDate(dateDoc)

	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return u, nil
}

// getItemById returns the document by id (completely)
func (d *Document) getItemById(id int64) (*DocItem, error) {
	di := DocItem{}
	sqlH := fmt.Sprintf("SELECT id, number, date, doc_type FROM %s WHERE id = $1", d.HeadersName)
	row := d.Db.QueryRow(sqlH, id)
	dateDoc := time.Time{}
	err := row.Scan(&di.Id, &di.Number, &dateDoc, &di.DocType)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	di.Number = di.GetNumber()
	di.Date = GetDate(dateDoc)

	sqlI := fmt.Sprintf("SELECT ri.row_id, ri.product_id, p.name, p.manufacturer_id, COALESCE(m.name, '') AS manufacturer_name, s.quantity, c.id AS cell_id, COALESCE(c.name, '') AS cell_name "+
		"FROM %s ri "+
		"LEFT JOIN %s p ON ri.product_id = p.id "+
		"LEFT JOIN %s m ON p.manufacturer_id = m.id "+
		"LEFT JOIN %s s ON s.doc_id = $1 AND s.row_id = ri.row_id "+
		"LEFT JOIN %s c ON s.cell_id = c.id "+
		"WHERE ri.parent_id = $1", d.ItemsName, tableRefProducts, tableRefManufacturers, "storage1", tableRefCells)
	rows, err := d.Db.Query(sqlI, id)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		r := DocRow{}
		cellId := 0
		cellName := ""
		err = rows.Scan(&r.RowId, &r.Product.Id, &r.Product.Name, &r.Product.Manufacturer.Id, &r.Product.Manufacturer.Name, &r.Quantity, &cellId, &cellName)
		if di.getType() == DocumentTypeReceipt {
			r.CellDst.Id = int64(cellId)
			r.CellDst.Name = cellName
		}
		di.Items = append(di.Items, r)
	}
	return &di, nil
}

func (d *Document) deleteItem(id int64) (int64, error) {
	sqlDelI := fmt.Sprintf("DELETE FROM %s WHERE parent_id = $1", d.ItemsName)
	sqlDelH := fmt.Sprintf("DELETE FROM %s WHERE id = $1", d.HeadersName)
	tx, err := d.Db.Begin()
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(sqlDelI, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	res, err := tx.Exec(sqlDelH, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return affRows, nil
}

func GetDate(date time.Time) string {
	return date.Format("02.01.2006 15:04:05")
}
