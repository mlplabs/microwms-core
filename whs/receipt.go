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
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlInsH := fmt.Sprintf("INSERT INTO %s (number, date, doc_type) VALUES($1, $2, $3) RETURNING id", tableDocReceiptHeaders)
	err = tx.QueryRow(sqlInsH, doc.Number, doc.Date, DocumentTypeReceipt).Scan(&doc.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	for idx, item := range doc.Items {

		pId, _, err := s.CreateProductInteractive(tx, item.Product.Name, item.Product.Manufacturer.Name, item.Product.ItemNumber, nil, nil)

		if err != nil {
			tx.Rollback()
			return 0, &core.WrapError{Err: err, Code: core.SystemError}
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
