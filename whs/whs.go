package whs

import (
	"fmt"
	"github.com/mlplabs/microwms-core/core"
)

// Whs is a physical warehouse object
// It must contain at least 3 Zone{} zones - acceptance, storage and shipment.
// There can be no more than 1 receiving and shipping zones, these zones are the entrance and exit in the warehouse, respectively
type Whs struct {
	Address        string `json:"address"`
	AcceptanceZone Zone   `json:"acceptance_zone"`
	ShippingZone   Zone   `json:"shipping_zone"`
	StorageZones   []Zone `json:"storage_zones"`
	CustomZones    []Zone `json:"custom_zones"`
	RefItem
}

// GetWhsItems returns a list of warehouses
func (s *Storage) GetWhsItems(offset int, limit int, parentId int64) ([]Whs, int, error) {
	items, count, err := s.GetReference(tableRefWhs).getItems(offset, limit, parentId)
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

// FindWhsById returns a warehouse item by id
func (s *Storage) FindWhsById(whsId int64) (*Whs, error) {
	item, err := s.GetReference(tableRefWhs).findItemById(whsId)
	u := new(Whs)
	u.RefItem = *item
	return u, err
}

// GetWhsById returns a warehouse object by id
func (s *Storage) GetWhsById(whsId int64) (*Whs, error) {
	w := new(Whs)
	w.StorageZones = make([]Zone, 0)
	w.CustomZones = make([]Zone, 0)

	sqlWhs := fmt.Sprintf("SELECT id, name, address FROM %s WHERE id = $1", tableRefWhs)
	row := s.Db.QueryRow(sqlWhs, whsId)

	err := row.Scan(&w.Id, &w.Name, &w.Address)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlZ := fmt.Sprintf("SELECT id, name, parent_id, zone_type FROM %s WHERE parent_id = $1", tableRefZones)
	rows, err := s.Db.Query(sqlZ, whsId)
	if err != nil {
		return nil, &core.WrapError{Err: err, Code: core.SystemError}
	}
	defer rows.Close()
	for rows.Next() {
		z := Zone{}
		err = rows.Scan(&z.Id, &z.Name, &z.ParentId, &z.ZoneType)
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

// GetWhsZones returns a warehouse object by id
func (s *Storage) GetWhsZones(whs *Whs) ([]Zone, error) {
	return s.FindZonesByParentId(whs.Id)
}

// CreateWhs creates a new warehouse
func (s *Storage) CreateWhs(u *Whs) (int64, error) {
	// TODO: необходимо создать storage

	tx, err := s.Db.Begin()
	if err != nil {
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	sqlCreate := fmt.Sprintf("INSERT INTO %s (name, address) VALUES ($1, $2) RETURNING id", "whs")
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

	sqlStorage := fmt.Sprintf(
		"create table if not exists storage%d ( "+
			"doc_id   integer default 0 not null, "+
			"doc_type smallint default 0 not null, "+
			"row_id   varchar(36) default ''::character varying not null, "+
			"row_time timestamptz default now() not null, "+
			"zone_id  integer, "+
			"cell_id  integer constraint storage%d_cells_id_fk references cells, "+
			"prod_id  integer,	"+
			"quantity integer ); "+
			"alter table storage%d owner to %s;", u.Id, u.Id, u.Id, s.dbUser)
	_, err = tx.Exec(sqlStorage)
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, &core.WrapError{Err: err, Code: core.SystemError}
	}

	return u.GetId(), nil
}

// UpdateWhs updates the warehouse
func (s *Storage) UpdateWhs(u *Whs) (int64, error) {
	return s.GetReference(tableRefWhs).updateItem(u)
}

// DeleteWhs deletes warehouse
func (s *Storage) DeleteWhs(u *Whs) (int64, error) {
	// TODO: need to remove child elements
	return s.GetReference(tableRefWhs).deleteItem(u.Id)
}
