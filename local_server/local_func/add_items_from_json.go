package localfunc

import (
	"database/sql"
	"fmt"
	"local_server/utils"
)

const (
	// Банчинг в turso
	RECORDS_PER_BATCH = 100 // 100 строк с 400 параметрами
	MAX_RETRIES       = 2   // 1 попытка + 2 ретрая
	BATCH_DELAY_MS    = 50  // Задержка между банчами
)

func AddItemsFromJson(db *sql.DB, path string) error {

	items, err := utils.ReadJson[[]Item](path)
	if err != nil {
		return err
	}

	query_prefix := "INSERT INTO items (id, name) VALUES "

	for i, item := range items {
		if i%100 == 0 {
			fmt.Printf("Processed %d items\n", i)
		}

		query := query_prefix + "(?, ?) ON CONFLICT(id) DO NOTHING"

		_, err := db.Exec(query, item.ID, item.Name)
		if err != nil {
			fmt.Printf("Error inserting item %d: %v\n", i, items[i])
			return err
		}
	}

	return nil
}
