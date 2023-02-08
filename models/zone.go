package models

// Типы зон
const (
	// ZoneTypeStorage - тип зоны: хранение
	ZoneTypeStorage = iota
	// ZoneTypeIncoming - тип зоны: приемка
	ZoneTypeIncoming
	// ZoneTypeOutGoing - тип зоны: отгрузка
	ZoneTypeOutGoing
	ZoneTypeCustom = 999
)

// Zone - зона склада
type Zone struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
	ZoneType int    `json:"zone_type"`
	RefItem
}

type ReferenceZones struct {
	Reference
}

// GetItems возвращает список зон
func (ref *ReferenceZones) GetItems(offset int, limit int, parentId int64) ([]Zone, int, error) {
	items, count, err := ref.getItems(offset, limit, parentId)
	if err != nil {
		return nil, 0, err
	}
	retVal := make([]Zone, len(items))
	for idx, item := range items {
		u := new(Zone)
		u.RefItem = item
		retVal[idx] = *u
	}
	return retVal, count, nil
}

// FindById выполняет поиск зоны по внутреннему идентификатору
func (ref *ReferenceZones) FindById(zoneId int64) (*Zone, error) {
	item, err := ref.findItemById(zoneId)
	u := new(Zone)
	u.RefItem = *item
	return u, err
}

// FindByParentId возвращает список зон для выбранного склада
func (ref *ReferenceZones) FindByParentId(whsId int64) ([]Zone, error) {
	sqlZones := "SELECT id, name FROM zones WHERE parent_id = $1"

	rows, err := ref.Db.Query(sqlZones, whsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Zone, 0)
	for rows.Next() {
		z := Zone{}
		err := rows.Scan(&z.Id, &z.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, z)
	}
	return res, nil
}
