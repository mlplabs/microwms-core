package models

// Manufacturer производитель
type Manufacturer struct {
	RefItem
}

type ReferenceManufacturers struct {
	Reference
	IRefSuggesting
}

// GetItems возвращает список производителей
func (ref *ReferenceManufacturers) GetItems(offset int, limit int, parentId int64) ([]Manufacturer, int, error) {
	items, count, err := ref.getItems(offset, limit, parentId)
	if err != nil {
		return nil, 0, err
	}
	retVal := make([]Manufacturer, len(items))

	for idx, item := range items {
		u := new(Manufacturer)
		u.RefItem = item
		retVal[idx] = *u
	}

	return retVal, count, nil
}

// FindById возвращает производителя по внутреннему идентификатору
func (ref *ReferenceManufacturers) FindById(mnfId int64) (*Manufacturer, error) {
	item, err := ref.findItemById(mnfId)
	u := new(Manufacturer)
	u.RefItem = *item
	return u, err
}

// FindByName возвращает список производителей по наименованию
func (ref *ReferenceManufacturers) FindByName(valName string) ([]Manufacturer, error) {
	items, err := ref.findItemByName(valName)
	if err != nil {
		return nil, err
	}
	retVal := make([]Manufacturer, len(items))
	for idx, item := range items {
		m := new(Manufacturer)
		m.RefItem = item
		retVal[idx] = *m
	}
	return retVal, err
}

func (ref *ReferenceManufacturers) GetSuggestion(text string, limit int) ([]Suggestion, error) {
	return ref.getSuggestion(text, limit)
}

// Create создает нового производителя
func (ref *ReferenceManufacturers) Create(m *Manufacturer) (int64, error) {
	return ref.createItem(m)
}

// Update обновляет производителя
func (ref *ReferenceManufacturers) Update(m *Manufacturer) (int64, error) {
	return ref.updateItem(m)
}

// Delete удаляет производителя
func (ref *ReferenceManufacturers) Delete(m *Manufacturer) (int64, error) {
	return ref.deleteItem(m.Id)

}
