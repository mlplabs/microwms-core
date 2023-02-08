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

type ReferenceBarcodes struct {
	Reference
}

// CreateBarcode создает новый штрих-код
func (ref *ReferenceBarcodes) CreateBarcode(b *Barcode) (int64, error) {
	sqlInsProd := fmt.Sprintf("INSERT INTO %s (barcode, barcode_type, parent_id) VALUES ($1, $2, $3) RETURNING id", ref.Name)
	err := ref.Db.QueryRow(sqlInsProd, b.Name, b.Type, b.ProdId).Scan(&b.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

// GetBarcodes возвращает список штрих-кодов
func (ref *ReferenceBarcodes) GetBarcodes(offset int, limit int) ([]Barcode, int, error) {

	var count int

	sqlBc := fmt.Sprintf("SELECT b.id, b.name, b.barcode_type, b.parent_id FROM %s b ORDER BY m.name ASC", ref.Name)

	if limit == 0 {
		limit = 10
	}
	rows, err := ref.Db.Query(sqlBc+" LIMIT $1 OFFSET $2", limit, offset)
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
	err = ref.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return bcs, count, nil
}

// FindBarcodeById возвращает штрих-код по внутреннему идентификатору
func (ref *ReferenceBarcodes) FindBarcodeById(bcId int64) (*Barcode, error) {
	sqlCell := fmt.Sprintf("SELECT b.id, b.name, b.type, b.parent_id FROM %s b WHERE b.id = $1", ref.Name)
	row := ref.Db.QueryRow(sqlCell, bcId)
	b := new(Barcode)
	err := row.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b, nil
}

// FindBarcodesByName возвращает список штрих-кодов по наименованию
// только по значению без типа и привязки
func (ref *ReferenceBarcodes) FindBarcodesByName(bcName string) ([]Barcode, error) {
	retBc := make([]Barcode, 0)

	sqlSel := fmt.Sprintf("SELECT id, name, type, parent_id FROM %s WHERE name = $1", ref.Name)
	rows, err := ref.Db.Query(sqlSel, bcName)
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
func (ref *ReferenceBarcodes) FindBarcodesByProdId(prodId int64) ([]Barcode, error) {
	fields := []string{"id", "name", "barcode_type", "parent_id"}
	fieldsStr := strings.Join(fields, ", ")

	pointers := make([]interface{}, len(fields))

	retBc := make([]Barcode, 0)
	sql := fmt.Sprintf("SELECT %s FROM %s WHERE name = $1", fieldsStr, ref.Parent)
	rows, err := ref.Db.Query(sql, prodId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		b := Barcode{}
		//err := rows.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
		err := rows.Scan(pointers...)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retBc = append(retBc, b)
	}
	return retBc, nil
}

// UpdateBarcode обновляет значение/тип/?привязку? штрих-кода
func (ref *ReferenceBarcodes) UpdateBarcode(b *Barcode) (int64, error) {
	sqlUpd := fmt.Sprintf("UPDATE %s SET barcode=$2, barcode_type=$3, parent_id=$4 WHERE id=$1", ref.Name)
	res, err := ref.Db.Exec(sqlUpd, b.Id, b.Name, b.Type, b.ProdId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

func (ref *ReferenceBarcodes) GetSuggestionBarcodes(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sqlSel := fmt.Sprintf("SELECT name FROM %s WHERE name LIKE $1 LIMIT $2", ref.Name)
	rows, err := ref.Db.Query(sqlSel, text+"%", limit)
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

func (ref *ReferenceBarcodes) DeleteBarcode(m *Barcode) (int64, error) {
	sqlDel := "DELETE FROM barcodes WHERE id=$1"
	res, err := ref.Db.Exec(sqlDel, m.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
