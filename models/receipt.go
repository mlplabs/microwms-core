package models

type DocumentReceipt struct {
	Document
}

func (dr *DocumentReceipt) GetItems(offset int, limit int) ([]DocItem, int, error) {
	return dr.getItems(offset, limit)
}

func (dr *DocumentReceipt) Create(doc *DocItem) (int64, error) {
	return dr.createItem(doc)
}

func (dr *DocumentReceipt) GetReceiptDocById(id int64) (*DocItem, error) {
	return dr.getItemById(id)
}
