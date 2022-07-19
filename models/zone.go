package models

// Типы зон
const (
	// ZoneTypeStorage - тип зоны: хранение
	ZoneTypeStorage = iota
	// ZoneTypeIncoming - тип зоны: приемка
	ZoneTypeIncoming
	// ZoneTypeOutGoing - тип зоны: отгрузка
	ZoneTypeOutGoing
	ZoneTypeCustom = 99
)

// Zone - зона склада
type Zone struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	WhsId    int    `json:"whs_id"`
	ZoneType int    `json:"zone_type"`
}

type ZoneService struct {
	Storage *Storage
}

// FindZoneById выполняет поиск зоны по внутреннему идентификатору
func (zs *ZoneService) FindZoneById(zoneId int64) (*Zone, error) {
	sqlCell := "SELECT id, name, whs_id, zone_type FROM zones WHERE id = $1"
	row := zs.Storage.Db.QueryRow(sqlCell, zoneId)
	z := new(Zone)
	err := row.Scan(&z.Id, &z.Name, &z.WhsId, &z.ZoneType)
	if err != nil {
		return nil, err
	}
	return z, nil
}

// GetZonesByWhsId возвращает список зон для выбранного склада
func (zs *ZoneService) GetZonesByWhsId(whsId int64) ([]Zone, error) {
	sqlZones := "SELECT id, name FROM zones WHERE whs_id = $1"

	rows, err := zs.Storage.Db.Query(sqlZones, whsId)
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
