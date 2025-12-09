package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

type Prices struct {
	Id         int      `json:"id"`
	Timestamp  string   `json:"timestamp"`
	IdAndPrice []string `json:"IdAndPrice"`
	Token      *string  `json:"tokens"`
}

func All_prices() ([]Prices, error) {
	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_TOKEN"))
	db, err := sql.Open("libsql", url)
	if err != nil {
		return []Prices{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM History_price_list")
	if err != nil {
		return []Prices{}, err
	}
	defer rows.Close()

	var prices []Prices
	for rows.Next() {
		var price Prices
		var idAndPriceRaw string // Временная переменная для сырой строки JSON

		err = rows.Scan(&price.Id, &idAndPriceRaw, &price.Timestamp, &price.Token)
		if err != nil {
			return []Prices{}, err
		}

		// Преобразуем JSON строку в массив строк
		if err := json.Unmarshal([]byte(idAndPriceRaw), &price.IdAndPrice); err != nil {
			return []Prices{}, fmt.Errorf("failed to parse IdAndPrice JSON: %v", err)
		}

		prices = append(prices, price)
	}

	return prices, nil
}
