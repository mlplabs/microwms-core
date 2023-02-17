package whs

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

// GetZonesItems returns a list of zones
func (s *Storage) GetZonesItems(offset int, limit int, parentId int64) ([]Zone, int, error) {
	items, count, err := s.GetReference(tableRefZones).getItems(offset, limit, parentId)
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

// FindZoneById searches for a zone by internal id
func (s *Storage) FindZoneById(zoneId int64) (*Zone, error) {
	item, err := s.GetReference(tableRefZones).findItemById(zoneId)
	u := new(Zone)
	u.RefItem = *item
	return u, err
}

// FindZonesByParentId returns a list of zones for the selected warehouse (by parent)
func (s *Storage) FindZonesByParentId(whsId int64) ([]Zone, error) {
	sqlZones := "SELECT id, name FROM zones WHERE parent_id = $1"
	rows, err := s.Db.Query(sqlZones, whsId)
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
