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
	Barcodes     []Barcode    `json:"barcodes"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Size         SpecificSize `json:"size"`
}

// Barcode объект штрих-кода
type Barcode struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}

type ProductService struct {
	Storage *Storage
}

// CreateProduct создает новый продукт
func (ps *ProductService) CreateProduct(p *Product) (int64, error) {
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
		mnfs, err := ps.FindManufacturerByName(p.Manufacturer.Name)
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

	tx, err := ps.Storage.Db.Begin()
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

	sqlInsProd := "INSERT INTO products (name, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
	err = tx.QueryRow(sqlInsProd, p.Name, mId, p.Size.Length, p.Size.Width, p.Size.Height, p.Size.Weight, p.Size.Volume, p.Size.UsefulVolume).Scan(&pId)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if p.Barcodes != nil {
		for _, bc := range p.Barcodes {
			sqlBc := "INSERT INTO barcodes (product_id, barcode, barcode_type) " +
				"VALUES($1, $2, $3) " +
				"ON CONFLICT (product_id, barcode, barcode_type) DO UPDATE SET product_id=$1, barcode=$2, barcode_type=$3"
			_, err := tx.Exec(sqlBc, pId, bc.Data, bc.Type)
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

// GetProductBarcodes возвращает список штрих-кодов продукта
func (ps *ProductService) GetProductBarcodes(productId int64) ([]Barcode, error) {
	var bcVal string
	var bcType int
	bcArr := make([]Barcode, 0, 0)

	sqlBc := "SELECT barcode, barcode_type FROM barcodes WHERE product_id = $1"
	rows, err := ps.Storage.Db.Query(sqlBc, productId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&bcVal, &bcType)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		b := Barcode{
			bcVal,
			bcType,
		}
		bcArr = append(bcArr, b)
	}

	return bcArr, nil
}

func (ps *ProductService) GetSuggestionProducts(text string, limit int) ([]string, error) {
	retVal := make([]string, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sql := "SELECT name FROM products WHERE name LIKE $1 LIMIT $2"
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

// FindProductById возвращает продукт по внутреннему идентификатору
func (ps *ProductService) FindProductById(productId int64) (*Product, error) {

	sqlCell := "SELECT p.id, p.name, p.manufacturer_id, m.name as manufacturer_name " +
		"FROM products p " +
		"LEFT JOIN manufacturers m ON p.manufacturer_id = m.id " +
		"WHERE p.id = $1"
	row := ps.Storage.Db.QueryRow(sqlCell, productId)
	p := new(Product)
	err := row.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	pBarcodes, err := ps.GetProductBarcodes(p.Id)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	p.Barcodes = pBarcodes
	return p, nil
}

// FindProductsByBarcode возвращает продукт по штрих-коду
func (ps *ProductService) FindProductsByBarcode(barcodeStr string) ([]Product, error) {
	var pId int64
	var bcType int
	var bcVal string
	prods := make([]Product, 0, 0)

	sqlBc := "SELECT product_id, barcode, barcode_type FROM barcodes WHERE barcode = $1"

	rows, err := ps.Storage.Db.Query(sqlBc, barcodeStr)
	if err != nil {
		return prods, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for rows.Next() {
		err := rows.Scan(&pId, &bcVal, &bcType)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p, err := ps.FindProductById(pId)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		prods = append(prods, *p)
	}

	return prods, nil
}

// GetProducts возвращает список продуктов
func (ps *ProductService) GetProducts(offset int, limit int) ([]Product, int, error) {
	var count int

	sqlProd := "SELECT p.id, p.name, p.manufacturer_id, m.name FROM products p " +
		"		LEFT JOIN manufacturers m ON p.manufacturer_id = m.id" +
		"		ORDER BY p.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := ps.Storage.Query(sqlProd+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	prods := make([]Product, count, 10)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)

		pBarcodes, err := ps.GetProductBarcodes(p.Id) // пока так
		if err != nil {
			return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p.Barcodes = pBarcodes

		prods = append(prods, *p)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlProd)
	err = ps.Storage.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return prods, count, nil
}

// UpdateProduct создает новый продукт
func (ps *ProductService) UpdateProduct(p *Product) (int64, error) {
	mId := int64(0)

	// сначала посмотрим производителя
	mId = p.Manufacturer.Id

	// если имя без id, то поищем сами и если будет 1, то его и возьмем
	if mId == 0 && strings.TrimSpace(p.Manufacturer.Name) != "" {
		mnfs, err := ps.FindManufacturerByName(p.Manufacturer.Name)
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

	tx, err := ps.Storage.Db.Begin()
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

	sqlInsProd := "UPDATE products SET name=$2,manufacturer_id=$3,sz_length=$4,sz_wight=$5,sz_height=$6,sz_weight=$7,sz_volume=$8, sz_uf_volume=$9 WHERE id=$1"

	res, err := tx.Exec(sqlInsProd, p.Id, p.Name, mId, p.Size.Length, p.Size.Width, p.Size.Height, p.Size.Weight, p.Size.Volume, p.Size.UsefulVolume)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if a, err := res.RowsAffected(); a != 1 || err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlDelBc := "DELETE FROM barcodes WHERE product_id=$1"
	res, err = tx.Exec(sqlDelBc, p.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if p.Barcodes != nil {
		for _, bc := range p.Barcodes {
			sqlBc := "INSERT INTO barcodes (product_id, barcode, barcode_type) " +
				"VALUES($1, $2, $3) " +
				"ON CONFLICT (product_id, barcode, barcode_type) DO UPDATE SET product_id=$1, barcode=$2, barcode_type=$3"
			_, err := tx.Exec(sqlBc, p.Id, bc.Data, bc.Type)
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

func (ps *ProductService) DeleteProduct(p *Product) (int64, error) {
	sqlDel := "DELETE FROM products WHERE id=$1"
	res, err := ps.Storage.Db.Exec(sqlDel, p.Id)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return affRows, nil
}
