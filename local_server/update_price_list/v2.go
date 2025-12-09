package updatepricelist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type UpdatePriceListV2 struct {
	Database *sql.DB
	Client   *http.Client
}

var Id_Markets = map[int]string{
	10000002: "Jita",
	//10000003: "Amarr",
	//10000004: "Dodixie",
	// 10000005: "Rens",
	// 10000006: "Hek",
}

var PriceList map[int]map[int]float64

const (
	//API ESI
	MAX_CONCURRENT_REQUESTS = 10 // Максимальное количество одновременных запросов

	// Банчинг в turso
	RECORDS_PER_BATCH = 100 // 100 строк с 400 параметрами
	MAX_RETRIES       = 2   // 1 попытка + 2 ретрая
	BATCH_DELAY_MS    = 50  // Задержка между банчами
)

func NewUpdatePriceListV2(db *sql.DB, client *http.Client) *UpdatePriceListV2 {

	return &UpdatePriceListV2{
		Database: db,
		Client:   client,
	}
}

func (u *UpdatePriceListV2) UpdatePriceList() error {
	item_ids, err := u.get_items_id()
	if err != nil {
		return err
	}

	PriceList = make(map[int]map[int]float64)

	for market_id := range Id_Markets {
		PriceList[market_id] = make(map[int]float64)
	}

	item_count := len(item_ids)

	semaphore := make(chan struct{}, MAX_CONCURRENT_REQUESTS)
	var wg sync.WaitGroup
	var mu sync.Mutex
	fmt.Printf("Начинается парсинг")
	for market_id := range Id_Markets {

		for i, id := range item_ids {
			wg.Add(1)

			go func(id, i, item_count, market_id int) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				price, err := u.get_price(id, i, item_count, market_id)
				if err != nil {
					fmt.Printf("Error getting price for item %d in market %d: %v\n", id, market_id, err)
					return
				}

				mu.Lock()
				PriceList[market_id][id] = price
				mu.Unlock()
			}(id, i, item_count, market_id)
		}
	}

	wg.Wait()
	close(semaphore)
	fmt.Printf("Парсинг завершен")

	// for market_id, prices := range PriceList {
	// 	fmt.Printf("Market %d prices: %v\n", market_id, prices)
	// 	for item_id, price := range prices {
	// 		fmt.Printf("Item %d price: %f\n", item_id, price)
	// 	}
	// }
	fmt.Printf("Начинается вставка данных")
	err = u.insertWithSmallBatches(PriceList)
	if err != nil {
		fmt.Printf("Error inserting prices: %v\n", err)
		return err
	}
	fmt.Printf("Вставка данных завершена")
	return nil
}
func (u *UpdatePriceListV2) insertWithSmallBatches(PriceList map[int]map[int]float64) error {

	type batchData struct {
		query  string
		values []any
	}

	var batches []batchData
	var currentValues []any
	var currentPlaceholders []string
	currentCount := 0

	for market_id, prices := range PriceList {
		for item_id, price := range prices {
			currentPlaceholders = append(currentPlaceholders, "(?, ?, ?, ?)")
			currentValues = append(currentValues, market_id, item_id, price, time.Now())
			currentCount++

			// Когда набрали batch
			if currentCount >= RECORDS_PER_BATCH {
				query := "INSERT INTO prices (market_id, item_id, price, timestamp) VALUES " +
					strings.Join(currentPlaceholders, ",")

				batches = append(batches, batchData{
					query:  query,
					values: currentValues,
				})

				// Сбрасываем для следующего банча
				currentPlaceholders = nil
				currentValues = nil
				currentCount = 0
			}
		}
	}

	// Добавляем остаток
	if currentCount > 0 {
		query := "INSERT INTO prices (market_id, item_id, price, timestamp) VALUES " +
			strings.Join(currentPlaceholders, ",")

		batches = append(batches, batchData{
			query:  query,
			values: currentValues,
		})
	}

	// Выполняем в транзакции
	tx, err := u.Database.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Вставляем каждый банч с ретраями
	for i, batch := range batches {
		var lastErr error

		for attempt := 0; attempt <= MAX_RETRIES; attempt++ {
			if attempt > 0 {
				// Экспоненциальная задержка
				delay := time.Duration(100*attempt) * time.Millisecond
				time.Sleep(delay)
				fmt.Printf("Retry %d for batch %d/%d\n", attempt, i+1, len(batches))
			}

			_, err = tx.Exec(batch.query, batch.values...)
			if err == nil {
				break // Успех
			}
			lastErr = err
		}

		if lastErr != nil {
			return fmt.Errorf("batch %d/%d failed after %d retries: %w",
				i+1, len(batches), MAX_RETRIES, lastErr)
		}

		// Пауза между банчами
		if i < len(batches)-1 {
			time.Sleep(time.Duration(BATCH_DELAY_MS) * time.Millisecond)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	totalRecords := 0
	for _, batch := range batches {
		totalRecords += len(batch.values) / 4
	}

	fmt.Printf("✅ Inserted %d records in %d batches\n", totalRecords, len(batches))
	return nil
}

func (u *UpdatePriceListV2) get_items_id() ([]int, error) {
	ids, err := u.Database.Query("SELECT id FROM items")
	if err != nil {
		return nil, err
	}
	defer ids.Close()
	var item_ids []int
	for ids.Next() {
		var id int
		if err := ids.Scan(&id); err != nil {
			return nil, err
		}
		item_ids = append(item_ids, id)
	}
	return item_ids, nil
}
func (u *UpdatePriceListV2) get_price(id int, i int, max int, market int) (float64, error) {
	// Очистка ID

	// Формирование URL
	url := fmt.Sprintf("https://esi.evetech.net/latest/markets/%d/orders/?type_id=%d&order_type=sell&language=en-us", market, id)

	resp, err := u.Client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %s", resp.Status)
	}

	// Чтение и парсинг JSON
	var orders []struct {
		Price      float64 `json:"price"`
		IsBuyOrder bool    `json:"is_buy_order"`
		// Можно добавить другие поля если нужны
	}

	// Используем json.Decoder для потокового чтения
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
		fmt.Println("Обработка", i+1, "из", max)
		//	fmt.Printf("Price for %s: %.2f\n", id, minPrice)
	}
	return minPrice, nil
}
