package models

// Whs объект физического склада
// Обязательно содержит минимум 3 зоны Zone{} - приемки, хранения и отгрузки.
// Зоны приемки и отгрузки не может быть более 1, эти зоны являются входом и выходом на складе соответственно
type Whs struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	AcceptanceZone Zone   `json:"-"`
	ShippingZone   Zone   `json:"-"`
	StorageZones   []Zone `json:"-"`
	RefItem
}

type ReferenceWarehouses struct {
	Reference
}

// GetItems возвращает список складов
func (ref *ReferenceWarehouses) GetItems(offset int, limit int, parentId int64) ([]Whs, int, error) {
	items, count, err := ref.getItems(offset, limit, parentId)
	if err != nil {
		return nil, 0, err
	}
	retVal := make([]Whs, len(items))
	for idx, item := range items {
		u := new(Whs)
		u.RefItem = item
		retVal[idx] = *u
	}
	return retVal, count, nil
}

// FindById возвращает склад по идентификатору
func (ref *ReferenceWarehouses) FindById(whsId int64) (*Whs, error) {
	item, err := ref.findItemById(whsId)
	u := new(Whs)
	u.RefItem = *item
	return u, err
}

// GetZones возвращает список зон склада
func (ref *ReferenceWarehouses) GetZones(whs *Whs) ([]Zone, error) {
	return ref.GetZonesByWhsId(whs.Id)
}

// GetZonesByWhsId	возвращает список зон склада по его идентификатору
func (ref *ReferenceWarehouses) GetZonesByWhsId(whsId int) ([]Zone, error) {
	sqlZones := "SELECT id, name, zone_type FROM zones WHERE parent_id = $1"

	rows, err := ref.Db.Query(sqlZones, whsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Zone, 0)
	for rows.Next() {
		z := Zone{}
		err := rows.Scan(&z.Id, &z.Name, &z.ZoneType)
		if err != nil {
			return nil, err
		}
		res = append(res, z)
	}
	return res, nil
}

// Create создает новый склад
func (ref *ReferenceWarehouses) Create(u *Whs) (int64, error) {
	// TODO: необходимо создать storage
	return ref.createItem(u)
}

// Update обновляет пользователя
func (ref *ReferenceWarehouses) Update(u *Whs) (int64, error) {
	return ref.updateItem(u)
}

// Delete удаляет склад
func (ref *ReferenceWarehouses) Delete(u *Whs) (int64, error) {
	// TODO: необходимо удаление дочерних элементов
	return ref.deleteItem(u.Id)
}
