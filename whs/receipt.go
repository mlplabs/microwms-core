package whs

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
)

func (s *Storage) GetReceiptDocsItems(offset int, limit int) ([]DocItem, int, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).getItems(offset, limit)
}

func (s *Storage) CreateReceiptDoc(doc *DocItem) (int64, error) {
	//return s.GetDocument(docTables{
	//	Headers: tableDocReceiptHeaders,
	//	Items:   tableDocReceiptItems,
	//}).createItem(doc)

	tx, err := s.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlInsH := fmt.Sprintf("INSERT INTO %s (number, date, doc_type) VALUES($1, $2, $3) RETURNING id", tableDocReceiptItems)
	err = tx.QueryRow(sqlInsH, doc.Number, doc.Date, DocumentTypeReceipt).Scan(&doc.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for idx, item := range doc.Items {
		mId := int64(0)
		pId := int64(0)

		mId = item.Product.Manufacturer.Id
		mnfs, err := s.FindManufacturersByName(item.Product.Manufacturer.Name)
		if err != nil {
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}
		if len(mnfs) > 1 {
			return 0, &core.WrapError{Err: fmt.Errorf("it is not possible to identify the manufacturer. found %d", len(mnfs)), Code: 0}
		}
		if len(mnfs) == 1 {
			mId = mnfs[0].Id
		}

		// If the manufacturer is not found by name, then we create it
		if mId == 0 {
			sqlIns := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", tableRefManufacturers)
			err = tx.QueryRow(sqlIns, item.Product.Manufacturer.Name).Scan(&mId)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}

		prods, err := s.FindProductsByName(item.Product.Name)
		if err != nil {
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
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
			sqlInsProd := fmt.Sprintf("INSERT INTO %s (name, item_number, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id", tableRefProducts)
			err = tx.QueryRow(sqlInsProd, item.Product.Name, item.Product.ItemNumber, mId, item.Product.Size.Length, item.Product.Size.Width, item.Product.Size.Height, item.Product.Size.Weight, item.Product.Size.Volume, item.Product.Size.UsefulVolume).Scan(&pId)
			if err != nil {
				tx.Rollback()
				return 0, &core.WrapError{Err: err, Code: core.SystemError}
			}
		}
		item.Product.Id = pId
		item.RowId = fmt.Sprintf("%d.%d", doc.Id, idx+1)
		sqlInsI := fmt.Sprintf("INSERT INTO %s (parent_id, row_id, product_id, quantity) VALUES($1, $2, $3, $4)", tableDocReceiptItems)
		_, err = tx.Exec(sqlInsI, doc.Id, item.RowId, item.Product.Id, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
		}

		c := Cell{Id: 2, WhsId: 1, ZoneId: 1}
		s := Storage{Db: s.Db}
		item.CellDst = c

		_, err = s.PutRow(doc, &item, tx)
		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}

		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return doc.Id, nil
}

func (s *Storage) GetReceiptDocById(id int64) (*DocItem, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).getItemById(id)
}

func (s *Storage) FindReceiptDocById(id int64) (*DocItem, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).findItemById(id)
}

func (s *Storage) UpdateReceiptDoc(doc *DocItem) (int64, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).updateItem(doc)
}

func (s *Storage) DeleteReceiptDoc(id int64) (int64, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).deleteItem(id)
}
