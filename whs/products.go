package whs

import (
	"database/sql"
	"fmt"
	"github.com/mlplabs/microwms-core/core"
	"strings"
)

// Product item, storage unit
type Product struct {
	ItemNumber   string       `json:"item_number"`
	Barcodes     []Barcode    `json:"barcodes"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Size         SpecificSize `json:"size"`
	RefItem
}

// GetProductsItems returns a list of products
func (s *Storage) GetProductsItems(offset int, limit int, parentId int64) ([]Product, int, error) {
	var count int

	sqlProd := "SELECT p.id, p.name, p.item_number, p.manufacturer_id, m.name As manufacturer_name FROM products p " +
		"		LEFT JOIN manufacturers m ON p.manufacturer_id = m.id" +
		"		ORDER BY p.name ASC"

	if limit == 0 {
		limit = 10
	}
	rows, err := s.Db.Query(sqlProd+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	prods := make([]Product, count)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.Id, &p.Name, &p.ItemNumber, &p.Manufacturer.Id, &p.Manufacturer.Name)

		pBarcodes, err := s.GetProductsBarcodes(p.Id) // пока так
		if err != nil {
			return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p.Barcodes = pBarcodes

		prods = append(prods, *p)
	}

	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM ( %s ) sub", sqlProd)
	err = s.Db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return nil, count, &core.WrapError{Err: err, Code: core.SystemError}
	}
	return prods, count, nil
}

// FindProductById returns product by internal id
func (s *Storage) FindProductById(productId int64) (*Product, error) {

	sqlCell := "SELECT p.id, p.name, p.item_number, p.manufacturer_id, m.name as manufacturer_name " +
		"FROM products p " +
		"LEFT JOIN manufacturers m ON p.manufacturer_id = m.id " +
		"WHERE p.id = $1"
	row := s.Db.QueryRow(sqlCell, productId)
	p := new(Product)
	err := row.Scan(&p.Id, &p.Name, &p.ItemNumber, &p.Manufacturer.Id, &p.Manufacturer.Name)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	pBarcodes, err := s.GetProductsBarcodes(p.Id)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	p.Barcodes = pBarcodes
	return p, nil
}

// FindProductsByName returns a list of products by name
func (s *Storage) FindProductsByName(valName string) ([]Product, error) {
	retItemList := make([]Product, 0)
	sql := fmt.Sprintf("SELECT id, name, manufacturer_id FROM %s WHERE name = $1", tableRefProducts)
	rows, err := s.Db.Query(sql, valName)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		item := Product{}
		err = rows.Scan(&item.Id, &item.Name, &item.Manufacturer.Id)
		if err != nil {
			return nil, &core.WrapError{Err: err, Code: core.SystemError}
		}
		retItemList = append(retItemList, item)
	}
	return retItemList, nil
}

// FindProductsByBarcode returns a product by barcode
func (s *Storage) FindProductsByBarcode(barcodeStr string) ([]Product, error) {
	var pId int64
	var bcType int
	var bcVal string
	prods := make([]Product, 0, 0)

	sqlBc := "SELECT parent_id, name, barcode_type FROM barcodes WHERE name = $1"

	rows, err := s.Db.Query(sqlBc, barcodeStr)
	if err != nil {
		return prods, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for rows.Next() {
		err := rows.Scan(&pId, &bcVal, &bcType)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		p, err := s.FindProductById(pId)
		if err != nil {
			return prods, &core.WrapError{Err: err, Code: core.SystemError}
		}
		prods = append(prods, *p)
	}

	return prods, nil
}

// GetProductsBarcodes returns a list of product barcodes
func (s *Storage) GetProductsBarcodes(productId int64) ([]Barcode, error) {

	var id int64
	var bcVal string
	var bcType int
	bcArr := make([]Barcode, 0, 0)

	sqlBc := "SELECT id, name, barcode_type FROM barcodes WHERE parent_id = $1"
	rows, err := s.Db.Query(sqlBc, productId)
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

func (s *Storage) GetProductsSuggestion(text string, limit int) ([]Suggestion, error) {
	retVal := make([]Suggestion, 0)

	if strings.TrimSpace(text) == "" {
		return retVal, &core.WrapError{Err: fmt.Errorf("invalid search text "), Code: core.SystemError}
	}
	if limit == 0 {
		limit = 10
	}

	sqlSel := fmt.Sprintf("SELECT p.id, p.name, m.name as mnf_name FROM %s p "+
		"LEFT JOIN %s m ON p.manufacturer_id = m.id "+
		"WHERE p.name ILIKE $1 LIMIT $2", tableRefProducts, tableRefManufacturers)
	rows, err := s.Db.Query(sqlSel, text+"%", limit)
	if err != nil {
		return retVal, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()

	for rows.Next() {
		item := Suggestion{}
		mnfName := ""
		err := rows.Scan(&item.Id, &item.Val, &mnfName)
		if err != nil {
			return retVal, &core.WrapError{Err: err, Code: core.SystemError}
		}
		item.Title = item.Val
		if mnfName != "" {
			item.Title += ", " + mnfName
		}
		retVal = append(retVal, item)
	}
	return retVal, err
}

// CreateProduct creates a new product
func (s *Storage) CreateProduct(p *Product) (int64, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	pId, _, err := s.CreateProductInteractive(tx, p.Name, p.Manufacturer.Name, p.ItemNumber, &p.Size, p.Barcodes)

	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return pId, nil
}

func (s *Storage) CreateProductInteractive(tx *sql.Tx, productName, manufacturerName, itemNumber string, size *SpecificSize, barcodes []Barcode) (int64, int64, error) {
	pId := int64(0)
	mId := int64(0)

	if strings.TrimSpace(productName) == "" {
		return 0, 0, &core.WrapError{Err: fmt.Errorf("product name is empty"), Code: 0}
	}

	// When creating a product, we always look for a manufacturer by name, regardless of its ID
	// because when creating, the user can choose from a list of hints and change the name
	mnfs, err := s.FindManufacturersByName(manufacturerName)
	if err != nil {
		return 0, 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if len(mnfs) > 1 {
		return 0, 0, &core.WrapError{Err: fmt.Errorf("it is not possible to identify the manufacturer. found %d", len(mnfs)), Code: 0}
	}
	if len(mnfs) == 1 {
		mId = mnfs[0].Id
	}
	// If the manufacturer is not found by name, then we create it
	if mId == 0 {
		sqlIns := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", tableRefManufacturers)
		err = tx.QueryRow(sqlIns, manufacturerName).Scan(&mId)
		if err != nil {
			tx.Rollback()
			return 0, 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
	}

	prods, err := s.FindProductsByName(productName)
	if err != nil {
		return 0, 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	// можем получить несколько товаров с одинаковым имененм
	// надо понять, имеется ли среди ниx товары с таким же производителем
	for _, v := range prods {
		if v.Manufacturer.Id == mId {
			// нашли существующий товар
			pId = v.Id
			break
		}
	}

	// если с таким именем и производителем не нашли, то создаем новый товар
	if pId == 0 {
		s := SpecificSize{}
		if size != nil {
			s = *size
		}

		sqlInsProd := fmt.Sprintf("INSERT INTO %s (name, item_number, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id", tableRefProducts)
		err = tx.QueryRow(sqlInsProd, productName, itemNumber, mId, s.Length, s.Width, s.Height, s.Weight, s.Volume, s.UsefulVolume).Scan(&pId)
		if err != nil {
			tx.Rollback()
			return 0, 0, &core.WrapError{Err: err, Code: core.SystemError}
		}

		if barcodes != nil {
			for _, bc := range barcodes {
				sqlBc := fmt.Sprintf("INSERT INTO %s (parent_id, name, barcode_type) "+
					"VALUES($1, $2, $3) "+
					"ON CONFLICT (parent_id, name, barcode_type) DO UPDATE SET parent_id=$1, name=$2, barcode_type=$3", tableRefBarcodes)
				_, err := tx.Exec(sqlBc, pId, bc.Name, bc.Type)
				if err != nil {
					tx.Rollback()
					return 0, 0, &core.WrapError{Err: err, Code: core.SystemError}
				}
			}
		}

	}
	return pId, mId, nil
}

// UpdateProduct updates the product
func (s *Storage) UpdateProduct(p *Product) (int64, error) {
	mId := int64(0)

	mnfs, err := s.FindManufacturersByName(p.Manufacturer.Name)
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	if len(mnfs) > 1 {
		return 0, &core.WrapError{Err: fmt.Errorf("it is not possible to identify the manufacturer. found %d", len(mnfs)), Code: core.SystemError}
	}
	if len(mnfs) == 1 {
		mId = mnfs[0].Id
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	if mId == 0 {
		sqlIns := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", tableRefManufacturers)
		err := tx.QueryRow(sqlIns, p.Manufacturer.Name).Scan(&mId)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
	}

	sqlInsProd := fmt.Sprintf("UPDATE %s SET name=$2,manufacturer_id=$3,sz_length=$4,sz_wight=$5,sz_height=$6,sz_weight=$7,sz_volume=$8, sz_uf_volume=$9, item_number=$10 WHERE id=$1", tableRefProducts)
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

// DeleteProduct removes the product
func (s *Storage) DeleteProduct(p *Product) (int64, error) {
	return s.GetReference(tableRefProducts).deleteItem(p.Id)
}
