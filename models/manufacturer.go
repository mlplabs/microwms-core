package models

// Manufacturer производитель
type Manufacturer struct {
	RefItem
}

type ReferenceManufacturers struct {
	Reference
}

// GetManufacturers возвращает список производителей
func (ref *ReferenceManufacturers) GetManufacturers(offset int, limit int) ([]Manufacturer, int, error) {
	items, count, err := ref.getItems(offset, limit)
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

// FindManufacturerById возвращает производителя по внутреннему идентификатору
func (ref *ReferenceManufacturers) FindManufacturerById(mnfId int64) (*Manufacturer, error) {
	item, err := ref.findItemById(mnfId)
	u := new(Manufacturer)
	u.RefItem = *item
	return u, err
}

// FindManufacturerByName возвращает список производителей по наименованию
func (ref *ReferenceManufacturers) FindManufacturerByName(valName string) ([]Manufacturer, error) {
	items, err := ref.findItemByName(valName)
	if err != nil {
		return nil, err
	}
	retVal := make([]Manufacturer, len(items))
	for _, item := range items {
		m := new(Manufacturer)
		m.RefItem = item
		retVal = append(retVal, *m)
	}
	return retVal, err
}

func (ref *ReferenceManufacturers) GetSuggestionManufacturers(text string, limit int) ([]string, error) {
	return ref.getSuggestion(text, limit)
}

// CreateManufacturer создает нового производителя
func (ref *ReferenceManufacturers) CreateManufacturer(m *Manufacturer) (int64, error) {
	return ref.createItem(m)
}

// UpdateManufacturer обновляет производителя
func (ref *ReferenceManufacturers) UpdateManufacturer(m *Manufacturer) (int64, error) {
	return ref.updateItem(m)
}

// DeleteManufacturer удаляет производителя
func (ref *ReferenceManufacturers) DeleteManufacturer(m *Manufacturer) (int64, error) {
	return ref.deleteItem(m.Id)

}
