package whs

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

// Barcode barcode object
type Barcode struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"` // barcode
	Type   int    `json:"type"`
	ProdId int64  `json:"prod_id"`
}

// CreateBarcode creates a new barcode
func (s *Storage) CreateBarcode(b *Barcode) (int64, error) {
	sqlInsProd := fmt.Sprintf("INSERT INTO %s (name, barcode_type, parent_id) VALUES ($1, $2, $3) RETURNING id", tableRefBarcodes)
	err := s.Db.QueryRow(sqlInsProd, b.Name, b.Type, b.ProdId).Scan(&b.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

// GetBarcodesItems returns a list of barcodes
func (s *Storage) GetBarcodesItems(offset int, limit int) ([]Barcode, int, error) {
	fields := []string{"id", "name", "barcode_type", "parent_id"}
	fieldsStr := strings.Join(fields, ", ")

	pointers := make([]interface{}, len(fields))

	var count int

	sqlBc := fmt.Sprintf("SELECT %s FROM %s ORDER BY name ASC", fieldsStr, tableRefBarcodes)

	if limit == 0 {
		limit = 10
	}
	rows, err := s.Db.Query(sqlBc+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	bcs := make([]Barcode, count, 10)
	for rows.Next() {
		b := new(Barcode)
		//err = rows.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
		err = rows.Scan(pointers...)
		bcs = append(bcs, *b)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlBc)
	err = s.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return bcs, count, nil
}

// FindBarcodeById returns barcode by internal ID
func (s *Storage) FindBarcodeById(bcId int64) (*Barcode, error) {
	sqlCell := fmt.Sprintf("SELECT b.id, b.name, b.barcode_type, b.parent_id FROM %s b WHERE b.id = $1", tableRefBarcodes)
	row := s.Db.QueryRow(sqlCell, bcId)
	b := new(Barcode)
	err := row.Scan(&b.Id, &b.Name, &b.Type, &b.ProdId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b, nil
}

// FindBarcodesByName returns a list of barcodes by name
// by value only without type and binding
func (s *Storage) FindBarcodesByName(bcName string) ([]Barcode, error) {
	retBc := make([]Barcode, 0)

	sqlSel := fmt.Sprintf("SELECT id, name, barcode_type, parent_id FROM %s WHERE name = $1", tableRefBarcodes)
	rows, err := s.Db.Query(sqlSel, bcName)
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

// FindBarcodesByProdId returns a list of barcodes for the product (owner)
func (s *Storage) FindBarcodesByProdId(prodId int64) ([]Barcode, error) {
	retBc := make([]Barcode, 0)
	sql := fmt.Sprintf("SELECT id, name, barcode_type, parent_id FROM %s WHERE name = $1", tableRefBarcodes)
	rows, err := s.Db.Query(sql, prodId)
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
func (s *Storage) UpdateBarcode(b *Barcode) (int64, error) {
	sqlUpd := fmt.Sprintf("UPDATE %s SET name=$2, barcode_type=$3, parent_id=$4 WHERE id=$1", tableRefBarcodes)
	res, err := s.Db.Exec(sqlUpd, b.Id, b.Name, b.Type, b.ProdId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if a, err := res.RowsAffected(); a != 1 || err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return b.Id, nil
}

func (s *Storage) GetSuggestionBarcodes(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sqlSel := fmt.Sprintf("SELECT name FROM %s WHERE name LIKE $1 LIMIT $2", tableRefBarcodes)
	rows, err := s.Db.Query(sqlSel, text+"%", limit)
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

func (s *Storage) DeleteBarcode(m *Barcode) (int64, error) {
	return s.DeleteBarcodeById(m.Id)
}

func (s *Storage) DeleteBarcodeById(barcodeId int64) (int64, error) {
	sqlDel := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableRefBarcodes)
	res, err := s.Db.Exec(sqlDel, barcodeId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
