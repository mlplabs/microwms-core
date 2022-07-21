package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
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

// Manufacturer производитель
type Manufacturer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductService struct {
	Storage *Storage
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
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for rows.Next() {
		err := rows.Scan(&pId, &bcVal, &bcType)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p, err := ps.FindProductById(pId)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		pBarcodes, err := ps.GetProductBarcodes(p.Id)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p.Barcodes = pBarcodes
		prods = append(prods, *p)
	}

	return prods, nil
}

// GetProducts возвращает список продуктов
func (ps *ProductService) GetProducts() ([]Product, error) {
	sqlProd := "SELECT p.id, p.name, p.manufacturer_id, m.name FROM products p LEFT JOIN manufacturers m ON p.manufacturer_id = m.id"
	rows, err := ps.Storage.Query(sqlProd)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	prods := make([]Product, 0, 10)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)

		pBarcodes, err := ps.GetProductBarcodes(p.Id) // пока так
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p.Barcodes = pBarcodes

		prods = append(prods, *p)
	}
	return prods, nil
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
	if p.Name == "" {
		return 0, &core.WrapError{Err: fmt.Errorf("required field 'name' is empty"), Code: 0}
	}

	mId = p.Manufacturer.Id

	if mId == 0 && p.Manufacturer.Name != "" {
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
		sqlIns := "INSERT INTO manufacturers (name) VALUES ($1)"
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

// GetManufacturers возвращает список производителей
func (ps *ProductService) GetManufacturers() ([]Manufacturer, error) {
	sqlMnf := "SELECT m.id, m.name FROM manufacturers m"
	rows, err := ps.Storage.Query(sqlMnf)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: 0}
	}
	defer rows.Close()

	mnfs := make([]Manufacturer, 0, 10)
	for rows.Next() {
		m := new(Manufacturer)
		err = rows.Scan(&m.Id, &m.Name)
		mnfs = append(mnfs, *m)
	}
	return mnfs, nil
}

// CreateManufacturer создает новый продукт
func (ps *ProductService) CreateManufacturer(m *Manufacturer) (int64, error) {
	mId := int64(0)
	if m.Id != 0 {
		return 0, &core.WrapError{Err: fmt.Errorf("possibly an error: create manufacturer with id <> 0"), Code: 0}
	}
	if m.Name == "" {
		return 0, &core.WrapError{Err: fmt.Errorf("required field 'name' is empty"), Code: 0}
	}

	sqlInsProd := "INSERT INTO manufacturers (name) VALUES ($1) RETURNING id"
	err := ps.Storage.Db.QueryRow(sqlInsProd, m.Name).Scan(&mId)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return mId, nil
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
