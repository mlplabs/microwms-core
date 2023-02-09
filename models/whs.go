package models

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
)

// Whs объект физического склада
// Обязательно содержит минимум 3 зоны Zone{} - приемки, хранения и отгрузки.
// Зоны приемки и отгрузки не может быть более 1, эти зоны являются входом и выходом на складе соответственно
type Whs struct {
	Address string `json:"address"`
	RefItem
}

type WhsObject struct {
	AcceptanceZone Zone   `json:"acceptance_zone"`
	ShippingZone   Zone   `json:"shipping_zone"`
	StorageZones   []Zone `json:"storage_zones"`
	CustomZones    []Zone `json:"custom_zones"`
	Whs
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

// FindById возвращает элемент склада по идентификатору
func (ref *ReferenceWarehouses) FindById(whsId int64) (*Whs, error) {
	item, err := ref.findItemById(whsId)
	u := new(Whs)
	u.RefItem = *item
	return u, err
}

// GetById возвращает объект склада по идентификатору
func (ref *ReferenceWarehouses) GetById(whsId int64) (*WhsObject, error) {
	w := new(WhsObject)
	w.StorageZones = make([]Zone, 0)
	w.CustomZones = make([]Zone, 0)

	sqlWhs := fmt.Sprintf("SELECT id, name, address FROM %s WHERE id = $1", ref.Name)
	row := ref.Db.QueryRow(sqlWhs, whsId)

	err := row.Scan(&w.Id, &w.Name, &w.Address)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlZ := "SELECT * FROM zones WHERE parent_id = $1"
	rows, err := ref.Db.Query(sqlZ, whsId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		z := Zone{}
		err = rows.Scan(&z.Id, &z.Name, &z.ZoneType)
		if z.ZoneType == ZoneTypeIncoming {
			w.AcceptanceZone = z
		}
		if z.ZoneType == ZoneTypeOutGoing {
			w.ShippingZone = z
		}
		if z.ZoneType == ZoneTypeStorage {
			w.StorageZones = append(w.StorageZones, z)
		}
		if z.ZoneType == ZoneTypeCustom {
			w.CustomZones = append(w.CustomZones, z)
		}

	}

	return w, nil
}

// GetZones возвращает список зон склада
func (ref *ReferenceWarehouses) GetZones(whs *Whs) ([]Zone, error) {
	return ref.GetZonesByWhsId(whs.Id)
}

// GetZonesByWhsId	возвращает список зон склада по его идентификатору
func (ref *ReferenceWarehouses) GetZonesByWhsId(whsId int64) ([]Zone, error) {
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

	tx, err := ref.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlCreate := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1, $2) RETURNING id", ref.Name)
	err = tx.QueryRow(sqlCreate, u.Name, u.Address).Scan(&u.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlCreateZIn := "INSERT INTO zones (name, zone_type, parent_id) VALUES ($1, $2, $3)"
	_, err = tx.Exec(sqlCreateZIn, "Зона приемки", ZoneTypeIncoming, u.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	_, err = tx.Exec(sqlCreateZIn, "Зона отгрузки", ZoneTypeOutGoing, u.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	_, err = tx.Exec(sqlCreateZIn, "Зона хранения", ZoneTypeStorage, u.Id)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}
	sqlStorage := fmt.Sprintf(
		"create table if not exists storage%d ( "+
			"zone_id  integer, "+
			"cell_id  integer constraint storage%d_cells_id_fk references cells, "+
			"prod_id  integer,	"+
			"quantity integer ); "+
			"alter table storage%d owner to devuser;", u.Id, u.Id, u.Id)
	ref.Db.Exec(sqlStorage)

	return u.GetId(), nil
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
