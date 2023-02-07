package models

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

// SpecificSize структура весогабаритных характеристик (см/см3/кг)
// полный объем: length * width * height
// полезный объем: length * width * height * K(0.8)
// вес: для продукта вес единицы в килограммах, для ячейки максимально возможный вес размещенных продуктов
type SpecificSize struct {
	Length       int     `json:"length"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Weight       float32 `json:"weight"`
	Volume       float32 `json:"volume"`
	UsefulVolume float32 `json:"useful_volume"` // Полезный объем ячейки
}

// Типы штрих-кодов
const (
	BarcodeTypeEAN13 = iota
	BarcodeTypeEAN8
	BarcodeTypeEAN14
	BarcodeTypeCode128
)

const (
	// CellDynamicPropIsService служебная ячейка. Запрещен автоматический отбор, но разрешены ручные перемещения в/из ячейки
	CellDynamicPropIsService = iota
	// CellDynamicPropSizeFree безразмерная ячейка. При размещении не используется проверка по размерам
	CellDynamicPropSizeFree
	// CellDynamicPropWeightFree при размещении не используется проверка по весу
	CellDynamicPropWeightFree
	// CellDynamicPropNotAllowedIn запрещено любое размещение в ячейку
	CellDynamicPropNotAllowedIn
	// CellDynamicPropNotAllowedOut запрещен любой отбор из ячейки
	CellDynamicPropNotAllowedOut
)

type Storage struct {
	Db *sql.DB
}

func (s *Storage) Init(host, dbname, dbuser, dbpass string) error {
	var err error
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, dbuser, dbpass)
	s.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = s.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// GetProductService возвращает менеджер для работы с продуктами
func (s *Storage) GetProductService() *ProductService {
	ps := new(ProductService)
	ps.Storage = s
	return ps
}

// GetWhsService возвращает менеджер для работы со складами
func (s *Storage) GetWhsService() *WhsService {
	ws := new(WhsService)
	ws.Storage = s
	return ws
}

// GetZoneService возвращает менеджер для работы с зонами склада
func (s *Storage) GetZoneService() *ZoneService {
	zs := new(ZoneService)
	zs.Storage = s
	return zs
}

// GetCellService возвращает менеджер для работы с ячейкам
func (s *Storage) GetCellService() *CellService {
	cs := new(CellService)
	cs.Storage = s
	return cs
}

func (s *Storage) GetUserService() *UserService {
	us := new(UserService)
	us.Storage = s
	us.Name = "users"
	return us
}

func (s *Storage) GetWarehouses() {

}

func (s *Storage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.Db.Query(query, args...)
}

// FindCellById возвращает ячейку по внутреннему идентификатору
func (s *Storage) FindCellById(cellId int64) (*Cell, error) {
	sqlCell := "SELECT id, name, whs_id, zone_id, passage_id, rack_id, floor, sz_length, sz_width, sz_height, sz_volume, sz_uf_volume, sz_weight, not_allowed_in, not_allowed_out, is_service, is_size_free, is_weight_free FROM cells WHERE id = $1"
	row := s.Db.QueryRow(sqlCell, cellId)
	c := new(Cell)

	err := row.Scan(&c.Id, &c.Name, &c.WhsId, &c.ZoneId, &c.PassageId, &c.RackId, &c.Floor,
		&c.Size.Length, &c.Size.Width, &c.Size.Height, &c.Size.Volume, &c.Size.UsefulVolume, &c.Size.Weight,
		&c.NotAllowedIn, &c.NotAllowedOut, &c.IsService, &c.IsSizeFree, &c.IsWeightFree)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Put размещает в ячейку (cell) продукт (prod) в количестве (quantity)
// Возвращает количество которое было размещено (quantity)
func (s *Storage) Put(cell *Cell, prod *Product, quantity int, tx *sql.Tx) (int, error) {
	var err error
	sqlIns := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", cell.WhsId)
	if tx != nil {
		_, err = tx.Exec(sqlIns, cell.ZoneId, cell.Id, prod.Id, quantity)
	} else {
		_, err = s.Db.Exec(sqlIns, cell.ZoneId, cell.Id, prod.Id, quantity)
	}
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

// Get отбирает из ячейки (cell) продукт (prod) в количестве (quantity)
// Возвращает отобранное количество (quantity)
func (s *Storage) Get(cell *Cell, prod *Product, quantity int, tx *sql.Tx) (int, error) {
	var err error

	if tx == nil {
		tx, err = s.Db.Begin()
		if err != nil {
			// не смогли начать транзакцию
			return 0, err
		}
	}

	sqlInsert := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", cell.WhsId)
	_, err = tx.Exec(sqlInsert, cell.ZoneId, cell.Id, prod.Id, -1*quantity)
	if err != nil {
		return 0, err
	}

	sqlQuant := fmt.Sprintf("SELECT SUM(quantity) AS quantity "+
		"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 AND prod_id = $3 "+
		"GROUP BY zone_id, cell_id, prod_id "+
		"HAVING SUM(quantity) < 0", cell.WhsId)
	rows, err := tx.Query(sqlQuant, cell.ZoneId, cell.Id, prod.Id)
	if err != nil {
		// ошибка контроля
		return 0, err
	}
	defer rows.Close()
	// мы должны получить пустой запрос
	if rows.Next() {
		err = tx.Rollback()
		if err != nil {
			// ошибка отката... все очень плохо
			return 0, err
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

// Quantity возвращает количество продуктов на св ячейке
func (s *Storage) Quantity(whsId int, cell Cell, tx *sql.Tx) (map[int]int, error) {
	var zoneId, cellId, prodId, quantity int
	res := make(map[int]int)

	sqlQuantity := fmt.Sprintf("SELECT zone_id, cell_id, prod_id, SUM(quantity) AS quantity "+
		"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 "+
		"GROUP BY zone_id, cell_id, prod_id "+
		"HAVING SUM(quantity) <> 0 %s", whsId, "")

	var err error
	var rows *sql.Rows

	if tx != nil {
		rows, err = tx.Query(sqlQuantity, cell.ZoneId, cell.Id)
	} else {
		rows, err = s.Db.Query(sqlQuantity, cell.ZoneId, cell.Id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&zoneId, &cellId, &prodId, &quantity)
		if err != nil {
			return nil, err
		}
		res[prodId] = quantity
	}
	return res, nil
}

func (s *Storage) Move(cellSrc, cellDst *Cell, prod *Product, quantity int) error {
	// TODO: cellSrc.WhsId <> cellDst.WhsId - временной разрыв или виртуальное перемещение

	_, err := s.Get(cellSrc, prod, quantity, nil)
	if err != nil {
		return err
	}
	_, err = s.Put(cellDst, prod, quantity, nil)
	if err == nil {
		return err
	}
	return nil
}

// BulkChangeSzCells устанавливает весогабаритные характеристики для массива ячеек
func (s *Storage) BulkChangeSzCells(cells []Cell, sz SpecificSize) (int64, error) {
	var ids []int64

	for _, c := range cells {
		ids = append(ids, c.Id)
	}
	sqlBulkUpdate := "UPDATE cells SET sz_length=$2, sz_width=$3, sz_height=$4, sz_volume=$5, sz_uf_volume=$6, sz_weight=$7 WHERE id = ANY($1)"
	res, err := s.Db.Exec(sqlBulkUpdate, pq.Array(ids), sz.Length, sz.Width, sz.Height, sz.Volume, sz.UsefulVolume, sz.Weight)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// BulkChangePropCells изменяет динамические параметры для массива ячеек
func (s *Storage) BulkChangePropCells(cells []Cell, CellDynamicProp int, value bool) (int64, error) {
	var ids []int64

	for _, c := range cells {
		ids = append(ids, c.Id)
	}

	cond := ""
	switch CellDynamicProp {
	case CellDynamicPropSizeFree:
		cond = "is_size_free = $1"
	case CellDynamicPropWeightFree:
		cond = "is_weight_free = $1"
	case CellDynamicPropNotAllowedIn:
		cond = "not_allowed_in = $1"
	case CellDynamicPropNotAllowedOut:
		cond = "not_allowed_out = $1"
	case CellDynamicPropIsService:
		cond = "is_service = $1"
	}

	sqlBulkUpdate := fmt.Sprintf("UPDATE %s WHERE id = ANY($1)", cond)
	res, err := s.Db.Exec(sqlBulkUpdate, pq.Array(ids), value)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
