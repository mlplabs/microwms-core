package models

// Manufacturer производитель
type Manufacturer struct {
	RefItem
}

// GetManufacturers возвращает список производителей
func (ps *ProductService) GetManufacturers(offset int, limit int) ([]Manufacturer, int, error) {
	items, count, err := ps.getItems(ps.Storage.Db, offset, limit)
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
func (ps *ProductService) FindManufacturerById(mnfId int64) (*Manufacturer, error) {
	item, err := ps.findItemById(ps.Storage.Db, mnfId)
	u := new(Manufacturer)
	u.RefItem = *item
	return u, err
}

// FindManufacturerByName возвращает список производителей по наименованию
func (ps *ProductService) FindManufacturerByName(valName string) ([]Manufacturer, error) {
	items, err := ps.findItemByName(ps.Storage.Db, valName)
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

func (ps *ProductService) GetSuggestionManufacturers(text string, limit int) ([]string, error) {
	return ps.getSuggestion(ps.Storage.Db, text, limit)
}

// CreateManufacturer создает нового производителя
func (ps *ProductService) CreateManufacturer(m *Manufacturer) (int64, error) {
	return ps.createItem(ps.Storage.Db, m)
}

// UpdateManufacturer обновляет производителя
func (ps *ProductService) UpdateManufacturer(m *Manufacturer) (int64, error) {
	return ps.updateItem(ps.Storage.Db, m)
}

// DeleteManufacturer удаляет производителя
func (ps *ProductService) DeleteManufacturer(m *Manufacturer) (int64, error) {
	return ps.deleteItem(ps.Storage.Db, m.Id)

}
