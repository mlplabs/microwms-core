package models

import "fmt"

// Cell ячейка склада
type Cell struct {
	Id            int64        `json:"id"`
	Name          string       `json:"name"`
	WhsId         int          `json:"whs_id"`     // Id склада (может быть именован)
	ZoneId        int          `json:"zone_id"`    // Id зоны назначения (может быть именован)
	PassageId     int          `json:"passage_id"` // Id проезда (может быть именован)
	RackId        int          `json:"rack_id"`    // Id стеллажа (может быть именован)
	Floor         int          `json:"floor"`
	IsSizeFree    bool         `json:"is_size_free"`
	IsWeightFree  bool         `json:"is_weight_free"`
	NotAllowedIn  bool         `json:"not_allowed_in"`
	NotAllowedOut bool         `json:"not_allowed_out"`
	IsService     bool         `json:"is_service"`
	Size          SpecificSize `json:"size"`
}

type CellService struct {
	Storage *Storage
}

// SetSize устанавливает размер ячейки
func (sz *SpecificSize) SetSize(length, width, height int, kUV float32) {
	sz.Volume = float32(length * width * height)
	sz.UsefulVolume = sz.Volume * kUV
}

// GetSize возвращает размеры ячейки
// length, width, height as int
// volume, usefulVolume as float
func (sz *SpecificSize) GetSize() (int, int, int, float32, float32) {
	return sz.Length, sz.Width, sz.Height, sz.Volume, sz.UsefulVolume
}

// GetNumeric возвращает строковое представление ячейки в виде набора чисел
func (c *Cell) GetNumeric() string {
	return fmt.Sprintf("%01d%02d%02d%02d%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}

// GetNumericView возвращает человеко-понятное представление (с разделителями)
func (c *Cell) GetNumericView() string {
	return fmt.Sprintf("%01d-%02d-%02d-%02d-%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}
