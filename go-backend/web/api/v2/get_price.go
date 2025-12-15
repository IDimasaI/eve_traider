package v2

import (
	"fmt"
)

func (api *Api2) Get_prices(id string) ([]PriceEntry, error) {
	query := `
        SELECT item_id, timestamp, price
        FROM prices
        WHERE item_id = ?
        ORDER BY timestamp DESC
    `
	rows, err := api.Database.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var prices []PriceEntry
	for rows.Next() {
		var p PriceEntry
		err := rows.Scan(&p.ItemID, &p.Timestamp, &p.Price)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}

		prices = append(prices, p)
	}

	return prices, nil
}
