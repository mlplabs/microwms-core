package whs

type RemainingProductRow struct {
	Product      RefItem `json:"product"`
	Manufacturer RefItem `json:"manufacturer"`
	Zone         RefItem `json:"zone"`
	Cell         Cell    `json:"cell"`
	Quantity     int     `json:"quantity"`
}

func (s *Storage) GetRemainingProducts() ([]RemainingProductRow, error) {
	retVal := make([]RemainingProductRow, 0)
	sqlSel := "SELECT store.prod_id AS product_id, coalesce(p.name, '<unnamed>') AS product_name, " +
		"       coalesce(m.id, 0) AS manufacturer_id, coalesce(m.name, '<unnamed>') AS manufacturer_name, " +
		"       store.zone_id, coalesce(z.name, '<unnamed>') AS zone_name, " +
		"       store.cell_id, c.name AS cell_name, " +
		"       store.quantity " +
		"FROM (SELECT s.prod_id, s.zone_id, s.cell_id, SUM(s.quantity) AS quantity " +
		"               FROM storage1 s " +
		"               GROUP BY s.prod_id, s.zone_id, s.cell_id) AS store " +
		"LEFT JOIN products p ON store.prod_id = p.id " +
		"LEFT JOIN manufacturers m on p.manufacturer_id = m.id " +
		"LEFT JOIN zones z ON store.zone_id = z.id " +
		"LEFT JOIN cells c ON store.cell_id = c.id " +
		"ORDER BY p.name"
	rows, err := s.Db.Query(sqlSel)
	if err != nil {
		return retVal, err
	}
	defer rows.Close()
	for rows.Next() {
		r := RemainingProductRow{}
		err = rows.Scan(&r.Product.Id, &r.Product.Name, &r.Manufacturer.Id, &r.Manufacturer.Name, &r.Zone.Id, &r.Zone.Name, &r.Cell.Id, &r.Cell.Name, &r.Quantity)
		if err != nil {
			return retVal, err
		}
		retVal = append(retVal, r)
	}

	return retVal, nil
}
