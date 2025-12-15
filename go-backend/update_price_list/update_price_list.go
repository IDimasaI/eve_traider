package updatepricelist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	localfunc "local_server/local_func"

	"net/http"
	"strings"
)

type UpdatePriceList struct {
	Database *sql.DB
	Client   *http.Client
}

func NewUpdatePriceList(db *sql.DB) *UpdatePriceList {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS History_price_list (id INTEGER PRIMARY KEY, IdAndPrice TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, tokens TEXT)")
	if err != nil {
		return nil
	}
	return &UpdatePriceList{
		Database: db,
		Client:   &http.Client{},
	}
}

type Items struct {
	items []string
}

// Ежечастно проверяем минимальный ордер на продажу(тот что люди продают)
func (u *UpdatePriceList) UpdatePriceList() error {
	var ids string
	err := u.Database.QueryRow("SELECT IdAndNames FROM all_ids").Scan(&ids)
	if err != nil {
		return err
	}

	var prices []string
	var id []localfunc.Items
	err = json.Unmarshal([]byte(ids), &id)
	if err != nil {
		return err
	}

	for i, _ := range id {

		price, err := u.get_price(id[i].Id, i, len(id))
		if err != nil {
			if strings.Contains(err.Error(), "no sell orders found for type_id") {
				fmt.Println(err)
			} else {
				return err
			}
		}

		prices = append(prices, fmt.Sprintf("%d:%.2f", id[i].Id, price))
	}
	//Для создания экранирующих "
	prices_text, _ := json.Marshal(prices)

	//utils.WriteFile("E:\\dev\\eve\\eve_traider\\local_server\\prices_db.txt", prices_text)
	//utils.WriteFile("E:\\dev\\eve\\eve_traider\\local_server\\ids_db.txt", []byte(ids))

	_, err = u.Database.Exec("INSERT INTO History_price_list (IdAndPrice) VALUES (?)", string(prices_text))
	if err != nil {
		return err
	}

	return nil
}

func (u *UpdatePriceList) get_price(id int, i int, max int) (float64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/markets/10000002/orders/?type_id=%d&order_type=sell&language=en-us", id)

	resp, err := u.Client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var orders []struct {
		Price      float64 `json:"price"`
		IsBuyOrder bool    `json:"is_buy_order"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return 0, fmt.Errorf("JSON decode failed: %w", err)
	}

	if len(orders) == 0 {
		return 0, nil
	}

	// Ищем минимальную цену среди ордеров на продажу (is_buy_order: false)
	var minPrice float64 = -1
	found := false

	for _, order := range orders {
		if !order.IsBuyOrder { // Ордера на продажу
			if !found || order.Price < minPrice {
				minPrice = order.Price
				found = true
			}
		}
	}

	if !found {
		return minPrice, fmt.Errorf("no sell orders found for type_id: %d", id)
	}
	// для каждых 10%
	if i%10 == 0 {
		fmt.Println("Обработка", i, "из", max)
		//	fmt.Printf("Price for %s: %.2f\n", id, minPrice)
	}
	return minPrice, nil
}
