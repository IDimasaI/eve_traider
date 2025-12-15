package localfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	Client *http.Client
}

type Item struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func NewClient(client *http.Client) *Client {
	return &Client{Client: client}
}

func (u *Client) CreateIDFromNames(names []string) ([]Item, error) {
	// Проверка входных данных
	if len(names) == 0 {
		return nil, fmt.Errorf("names slice is empty")
	}

	// 1. Формируем правильный JSON
	jsonData, err := json.Marshal(names)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 2. Создаем запрос с правильными заголовками
	req, err := http.NewRequest(
		"POST",
		"https://esi.evetech.net/universe/ids",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Устанавливаем обязательные заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Compatibility-Date", "2025-11-06") // Актуальная дата

	// 4. Отправляем запрос
	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// 5. Проверяем статус (ESI может возвращать разные успешные коды)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// ИЛИ более универсальная структура для поиска
	var esiResponse map[string][]struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esiResponse); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}

	// 7. Сопоставляем имена с ID
	result := make([]Item, len(names))
	foundCount := 0

	// Ищем во всех категориях ответа
	for name_cat, categoryItems := range esiResponse {
		if name_cat != "inventory_types" {
			continue
		}
		fmt.Println("Категория:", name_cat)
		for _, item := range categoryItems {
			for i, name := range names {
				if item.Name == name {
					result[i] = Item{ID: item.ID, Name: item.Name}
					foundCount++
					break
				}
			}
		}
	}

	if foundCount == 0 {
		return nil, fmt.Errorf("no IDs found for given names")
	}

	return result, nil
}

// JSONArrayExtractor - основной инструмент для извлечения массивов
type JSONArrayExtractor struct {
	// Сохранять информацию о пути
	KeepPathInfo bool
	// Только уникальные значения
	UniqueOnly bool
	// Фильтр по минимальному размеру массива
	MinArraySize int
}

func NewExtractor() *JSONArrayExtractor {
	return &JSONArrayExtractor{
		KeepPathInfo: false,
		UniqueOnly:   true,
		MinArraySize: 0,
	}
}

// Extract извлекает массивы из JSON
func (e *JSONArrayExtractor) Extract(jsonData []byte) ([]string, map[string][]string, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, nil, err
	}

	flatList := []string{}
	grouped := make(map[string][]string)
	seen := make(map[string]bool)

	var recursiveExtract func(interface{}, string)
	recursiveExtract = func(node interface{}, path string) {
		switch v := node.(type) {
		case map[string]interface{}:
			for key, val := range v {
				newPath := key
				if path != "" {
					newPath = path + "." + key
				}
				recursiveExtract(val, newPath)
			}

		case []interface{}:
			// Проверяем минимальный размер
			if e.MinArraySize > 0 && len(v) < e.MinArraySize {
				return
			}

			// Извлекаем строки из массива
			var stringsInArray []string
			for _, item := range v {
				if str, ok := item.(string); ok {
					stringsInArray = append(stringsInArray, str)
				}
			}

			// Если массив содержит строки
			if len(stringsInArray) > 0 {
				if e.KeepPathInfo {
					grouped[path] = stringsInArray
				}

				// Добавляем в плоский список
				for _, str := range stringsInArray {
					if e.UniqueOnly {
						if !seen[str] {
							seen[str] = true
							flatList = append(flatList, str)
						}
					} else {
						flatList = append(flatList, str)
					}
				}
			} else {
				// Рекурсивно обходим нестроковые элементы массива
				for i, item := range v {
					newPath := fmt.Sprintf("%s[%d]", path, i)
					recursiveExtract(item, newPath)
				}
			}
		}
	}

	recursiveExtract(data, "")
	return flatList, grouped, nil
}

func Test_storage(all_names []string) {
	total := len(all_names)

	fmt.Println("count item: ", total)

	if total > 500 {
		fmt.Println("Борщь пи%*;%")
	}

	count_update_in_day := 4
	updates_per_day := 24 / count_update_in_day

	fmt.Printf("За 24ч строк: %d при обновлении %dч \n",
		total*updates_per_day,
		count_update_in_day)

	fmt.Printf("За 30д строк: %d \n",
		total*updates_per_day*30)

	estimateStorage(total, count_update_in_day, 1, 30)
}

func estimateStorage(itemsCount int, updateIntervalHours int, marketsCount int, days int) {
	updatesPerDay := 24 / updateIntervalHours
	rowSize := 72 // байт
	totalRows := itemsCount * updatesPerDay * marketsCount * days
	storageMB := float64(totalRows*rowSize) / (1024 * 1024)

	fmt.Printf("Оценка хранилища для %d предметов, обновление каждые %d часа, %d рынков, за %d дней:\n",
		itemsCount, updateIntervalHours, marketsCount, days)
	fmt.Printf("  Всего записей: %d\n", totalRows)
	fmt.Printf("  Объем: %.2f МБ\n", storageMB)
}

func (client *Client) Chunk_worker(all_names []string, chunkSize int) ([]Item, error) {
	var ids []Item
	var err error
	if len(all_names) > chunkSize {
		for i := 0; i < len(all_names); i += chunkSize {
			end := i + chunkSize
			if end > len(all_names) {
				end = len(all_names)
			}
			subset := all_names[i:end]
			subset_ids, err := client.CreateIDFromNames(subset)
			fmt.Println("Обработано:", len(subset_ids))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create ids: %s", err)
				return nil, err
			}
			ids = append(ids, subset_ids...)
			time.Sleep(1 * time.Second)
		}
	} else {
		ids, err = client.CreateIDFromNames(all_names)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
