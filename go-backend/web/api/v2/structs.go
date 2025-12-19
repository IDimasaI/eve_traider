package v2

import (
	"database/sql"
	"time"
)

type ProxyCache struct {
	//TODO: добавить кеширование запросов.
}

type Api2 struct {
	Database   *sql.DB
	ProxyCache ProxyCache
}

func NewApi2(db *sql.DB) *Api2 {
	return &Api2{
		Database: db,
	}
}

type Item struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Category       *string `json:"category"`
	Is_observation bool    `json:"is_observation"`
}

type PriceEntry struct {
	ItemID    int       `json:"item_id"`
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
	Market_id int       `json:"market_id"`
}
