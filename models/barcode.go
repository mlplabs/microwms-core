package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

// Barcode объект штрих-кода
type Barcode struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"` // barcode
	Type   int    `json:"type"`
	ProdId int64  `json:"prod_id"`
}

// CreateBarcode создает новый штрих-код
func (ps *ProductService) CreateBarcode(b *Barcode) (int64, error) {
	sqlInsProd := "INSERT INTO barcodes (barcode, barcode_type, product_id) VALUES ($1, $2, $3) RETURNING id"
	err := ps.Storage.Db.QueryRow(sqlInsProd, b.Name, b.Type, b.ProdId).Scan(&b.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

// GetBarcodes возвращает список штрих-кодов
func (ps *ProductService) GetBarcodes(offset int, limit int) ([]Barcode, int, error) {
	var count int

	sqlBc := "SELECT b.id, b.name, b.barcode_type, b.product_id FROM barcodes b " +
		"		ORDER BY m.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := ps.Storage.Query(sqlBc+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	bcs := make([]Barcode, count, 10)
	for rows.Next() {
		b := new(Barcode)
		err = rows.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
		bcs = append(bcs, *b)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlBc)
	err = ps.Storage.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return bcs, count, nil
}

// FindBarcodeById возвращает штрих-код по внутреннему идентификатору
func (ps *ProductService) FindBarcodeById(bcId int64) (*Barcode, error) {
	sqlCell := "SELECT b.id, b.name, b.type, b.product_id " +
		"FROM barcodes b " +
		"WHERE b.id = $1"
	row := ps.Storage.Db.QueryRow(sqlCell, bcId)
	b := new(Barcode)
	err := row.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b, nil
}

// FindBarcodesByName возвращает список штрих-кодов по наименованию
// только по значению без типа и привязки
func (ps *ProductService) FindBarcodesByName(bcName string) ([]Barcode, error) {
	retBc := make([]Barcode, 0)
	sql := "SELECT id, name, type, product_id FROM barcodes WHERE name = $1"
	rows, err := ps.Storage.Query(sql, bcName)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		b := Barcode{}
		err := rows.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retBc = append(retBc, b)
	}
	return retBc, nil
}

// FindBarcodesByProdId возвращает список штрих-кодов по товару (владельцу)
func (ps *ProductService) FindBarcodesByProdId(prodId int64) ([]Barcode, error) {
	retBc := make([]Barcode, 0)
	sql := "SELECT id, name, type, product_id FROM products WHERE name = $1"
	rows, err := ps.Storage.Query(sql, prodId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		b := Barcode{}
		err := rows.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retBc = append(retBc, b)
	}
	return retBc, nil
}

// UpdateBarcode обновляет значение/тип/?привязку? штрих-кода
func (ps *ProductService) UpdateBarcode(b *Barcode) (int64, error) {
	sqlUpd := "UPDATE barcodes SET barcode=$2, barcode_type=$3, product_id=$4 WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlUpd, b.Id, b.Name, b.Type, b.ProdId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

func (ps *ProductService) GetSuggestionBarcodes(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sql := "SELECT name FROM barcodes WHERE name LIKE $1 LIMIT $2"
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

func (ps *ProductService) DeleteBarcode(m *Barcode) (int64, error) {
	sqlDel := "DELETE FROM barcodes WHERE id=$1"
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
