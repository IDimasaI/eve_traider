package v2

import (
	"fmt"
)

func (api *Api2) Get_all_items() ([]Item, error) {

	query := `
        SELECT id, name, category, is_observation
        FROM items
    `

	rows, err := api.Database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var p Item
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Is_observation)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, p)
	}

	return items, nil
}
