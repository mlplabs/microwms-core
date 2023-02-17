package whs

func (s *Storage) GetReceiptDocsItems(offset int, limit int) ([]DocItem, int, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).getItems(offset, limit)
}

func (s *Storage) CreateReceiptDoc(doc *DocItem) (int64, error) {
	return s.GetDocument(docTables{
		Headers: tableDocReceiptHeaders,
		Items:   tableDocReceiptItems,
	}).createItem(doc)
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
