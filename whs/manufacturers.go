package whs

// Manufacturer manufacturer
type Manufacturer struct {
	RefItem
}

// GetManufacturersItems returns a list of manufacturers
func (s *Storage) GetManufacturersItems(offset int, limit int, parentId int64) ([]Manufacturer, int, error) {
	items, count, err := s.GetReference(tableRefManufacturers).getItems(offset, limit, parentId)
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

// FindManufacturerById returns manufacturer by internal id
func (s *Storage) FindManufacturerById(mnfId int64) (*Manufacturer, error) {
	item, err := s.GetReference(tableRefManufacturers).findItemById(mnfId)
	u := new(Manufacturer)
	u.RefItem = *item
	return u, err
}

// FindManufacturersByName returns a list of manufacturers by name
func (s *Storage) FindManufacturersByName(valName string) ([]Manufacturer, error) {
	items, err := s.GetReference(tableRefManufacturers).findItemByName(valName)
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

func (s *Storage) GetManufacturersSuggestion(text string, limit int) ([]Suggestion, error) {
	return s.GetReference(tableRefManufacturers).getSuggestion(text, limit)
}

// CreateManufacturer creates a new producer
func (s *Storage) CreateManufacturer(m *Manufacturer) (int64, error) {
	return s.GetReference(tableRefManufacturers).createItem(m)
}

// UpdateManufacturer updates the manufacturer
func (s *Storage) UpdateManufacturer(m *Manufacturer) (int64, error) {
	return s.GetReference(tableRefManufacturers).updateItem(m)
}

// DeleteManufacturer removes the manufacturer
func (s *Storage) DeleteManufacturer(m *Manufacturer) (int64, error) {
	return s.GetReference(tableRefManufacturers).deleteItem(m.Id)
}
