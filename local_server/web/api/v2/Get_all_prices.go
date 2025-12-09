package v2

import (
	"fmt"
)

func (api *Api2) Get_all_prices() ([]PriceEntry, error) {
	query := `
        SELECT item_id, timestamp, price
        FROM prices

        ORDER BY timestamp DESC, item_id
    `

	rows, err := api.Database.Query(query)
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
