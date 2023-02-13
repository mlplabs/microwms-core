package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

// Product единица хранения
type Product struct {
	Id           int64        `json:"id"`
	Name         string       `json:"name"`
	ItemNumber   string       `json:"item_number"`
	Barcodes     []Barcode    `json:"barcodes"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Size         SpecificSize `json:"size"`
}

type ReferenceProducts struct {
	Reference
}

// GetItems возвращает список продуктов
func (ref *ReferenceProducts) GetItems(offset int, limit int, parentId int64) ([]Product, int, error) {
	var count int

	sqlProd := "SELECT p.id, p.name, p.item_number, p.manufacturer_id, m.name FROM products p " +
		"		LEFT JOIN manufacturers m ON p.manufacturer_id = m.id" +
		"		ORDER BY p.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := ref.Db.Query(sqlProd+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	prods := make([]Product, count, 10)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.Id, &p.Name, &p.ItemNumber, &p.Manufacturer.Id, &p.Manufacturer.Name)

		pBarcodes, err := ref.GetBarcodes(p.Id) // пока так
		if err != nil {
			return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p.Barcodes = pBarcodes

		prods = append(prods, *p)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlProd)
	err = ref.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return prods, count, nil
}

// FindById возвращает продукт по внутреннему идентификатору
func (ref *ReferenceProducts) FindById(productId int64) (*Product, error) {

	sqlCell := "SELECT p.id, p.name, p.item_number, p.manufacturer_id, m.name as manufacturer_name " +
		"FROM products p " +
		"LEFT JOIN manufacturers m ON p.manufacturer_id = m.id " +
		"WHERE p.id = $1"
	row := ref.Db.QueryRow(sqlCell, productId)
	p := new(Product)
	err := row.Scan(&p.Id, &p.Name, &p.ItemNumber, &p.Manufacturer.Id, &p.Manufacturer.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	pBarcodes, err := ref.GetBarcodes(p.Id)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	p.Barcodes = pBarcodes
	return p, nil
}

// FindByBarcode возвращает продукт по штрих-коду
func (ref *ReferenceProducts) FindByBarcode(barcodeStr string) ([]Product, error) {
	var pId int64
	var bcType int
	var bcVal string
	prods := make([]Product, 0, 0)

	sqlBc := "SELECT parent_id, name, barcode_type FROM barcodes WHERE name = $1"

	rows, err := ref.Db.Query(sqlBc, barcodeStr)
	if err != nil {
		return prods, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for rows.Next() {
		err := rows.Scan(&pId, &bcVal, &bcType)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p, err := ref.FindById(pId)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		prods = append(prods, *p)
	}

	return prods, nil
}

// GetBarcodes возвращает список штрих-кодов продукта
func (ref *ReferenceProducts) GetBarcodes(productId int64) ([]Barcode, error) {

	var id int64
	var bcVal string
	var bcType int
	bcArr := make([]Barcode, 0, 0)

	sqlBc := "SELECT id, name, barcode_type FROM barcodes WHERE parent_id = $1"
	rows, err := ref.Db.Query(sqlBc, productId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &bcVal, &bcType)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		b := Barcode{
			id,
			bcVal,
			bcType,
			productId,
		}
		bcArr = append(bcArr, b)
	}

	return bcArr, nil
}

func (ref *ReferenceProducts) FindManufacturerByName(mnfName string) ([]Manufacturer, error) {
	m := &ReferenceManufacturers{
		Reference: Reference{
			Name: "manufacturers",
			Db:   ref.Db,
		},
	}
	return m.FindByName(mnfName)
}

func (ref *ReferenceProducts) GetSuggestion(text string, limit int) ([]Suggestion, error) {
	return ref.getSuggestion(text, limit)
}

// Create создает новый продукт
func (ref *ReferenceProducts) Create(p *Product) (int64, error) {
	pId := int64(0)
	mId := int64(0)

	if p.Id != 0 {
		return 0, &core.WrapError{
			Err:  fmt.Errorf("possibly an error: create product with id <> 0"),
			Code: 0,
		}
	}
	if strings.TrimSpace(p.Name) == "" {
		return 0, &core.WrapError{Err: fmt.Errorf("required field 'name' is empty"), Code: 0}
	}

	mId = p.Manufacturer.Id

	if mId == 0 && strings.TrimSpace(p.Manufacturer.Name) != "" {
		mnfs, err := ref.FindManufacturerByName(p.Manufacturer.Name)
		if err != nil {
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
		if len(mnfs) > 1 {
			return 0, &core.WrapError{Err: fmt.Errorf("it is not possible to identify the manufacturer. found %d", len(mnfs)), Code: 0}
		}
		if len(mnfs) == 1 {
			mId = mnfs[0].Id
		}
	}

	tx, err := ref.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if mId == 0 {
		sqlIns := "INSERT INTO manufacturers (name) VALUES ($1) RETURNING id"
		err := tx.QueryRow(sqlIns, p.Manufacturer.Name).Scan(&mId)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
	}

	sqlInsProd := "INSERT INTO products (name, item_number, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"
	err = tx.QueryRow(sqlInsProd, p.Name, p.ItemNumber, mId, p.Size.Length, p.Size.Width, p.Size.Height, p.Size.Weight, p.Size.Volume, p.Size.UsefulVolume).Scan(&pId)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if p.Barcodes != nil {
		for _, bc := range p.Barcodes {
			sqlBc := "INSERT INTO barcodes (parent_id, name, barcode_type) " +
				"VALUES($1, $2, $3) " +
				"ON CONFLICT (parent_id, name, barcode_type) DO UPDATE SET parent_id=$1, name=$2, barcode_type=$3"
			_, err := tx.Exec(sqlBc, pId, bc.Name, bc.Type)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return pId, nil
}

// Update создает новый продукт
func (ref *ReferenceProducts) Update(p *Product) (int64, error) {
	mId := int64(0)

	// сначала посмотрим производителя
	//	mId = p.Manufacturer.Id // Производителя могли поменять - ищем по имени или создаем

	if mId == 0 && strings.TrimSpace(p.Manufacturer.Name) != "" {
		mnfs, err := ref.FindManufacturerByName(p.Manufacturer.Name)
		if err != nil {
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
		if len(mnfs) > 1 {
			return 0, &core.WrapError{Err: fmt.Errorf("it is not possible to identify the manufacturer. found %d", len(mnfs)), Code: core.SystemError}
		}
		if len(mnfs) == 1 {
			mId = mnfs[0].Id
		}
	}

	tx, err := ref.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if mId == 0 {
		sqlIns := "INSERT INTO manufacturers (name) VALUES ($1) RETURNING id"
		err := tx.QueryRow(sqlIns, p.Manufacturer.Name).Scan(&mId)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
	}

	sqlInsProd := "UPDATE products SET name=$2,manufacturer_id=$3,sz_length=$4,sz_wight=$5,sz_height=$6,sz_weight=$7,sz_volume=$8, sz_uf_volume=$9, item_number=$10 WHERE id=$1"

	res, err := tx.Exec(sqlInsProd, p.Id, p.Name, mId, p.Size.Length, p.Size.Width, p.Size.Height, p.Size.Weight, p.Size.Volume, p.Size.UsefulVolume, p.ItemNumber)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if a, err := res.RowsAffected(); a != 1 || err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlDelBc := "DELETE FROM barcodes WHERE parent_id=$1"
	res, err = tx.Exec(sqlDelBc, p.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if p.Barcodes != nil {
		for _, bc := range p.Barcodes {
			sqlBc := "INSERT INTO barcodes (parent_id, name, barcode_type) " +
				"VALUES($1, $2, $3) " +
				"ON CONFLICT (parent_id, name, barcode_type) DO UPDATE SET parent_id=$1, name=$2, barcode_type=$3"
			_, err := tx.Exec(sqlBc, p.Id, bc.Name, bc.Type)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return p.Id, nil
}

// Delete удаляет продукт
func (ref *ReferenceProducts) Delete(p *Product) (int64, error) {
	return ref.deleteItem(p.Id)
}
