package run

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"local_server/utils"
	"strings"
)

type Migration_str struct {
	Id        int    `json:"id"`
	Ids       string `json:"ids"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
}
type Migration_new struct {
	Id         int    `json:"id"`
	IdAndPrice string `json:"id_and_price"`
	Timestamp  string `json:"timestamp"`
}

func Migration(db *sql.DB) error {
	if err := createTable_if_not(db); err != nil {
		return err
	}

	file, err := utils.ReadJson[[]Migration_str]("C:\\Users\\Dmitriy\\Downloads\\History_price_list.json")
	if err != nil {
		return err
	}

	new_file := make([]Migration_new, len(file))
	for i, migration := range file {
		new_file[i] = generated_new(migration)
	}

	for _, migration := range new_file {
		_, err := db.Exec("INSERT INTO History_price_list (IdAndPrice, timestamp) VALUES (?, ?)", migration.IdAndPrice, migration.Timestamp)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

type IdAndPrice struct {
	Id    string `json:"id"`
	Price string `json:"price"`
}

func generated_new(data Migration_str) Migration_new {
	var new Migration_new

	prices := strings.Split(strings.ReplaceAll(data.Price, "\"", ""), ",")
	ids := strings.Split(strings.ReplaceAll(data.Ids, "\"", ""), ",")

	new_price_and_id := make([]string, len(ids))
	for i, id := range ids {

		price := prices[i]
		new_price_and_id[i] = fmt.Sprintf("%s:%s", id, price)
	}

	text, _ := json.Marshal(new_price_and_id)

	new.Id = data.Id
	new.IdAndPrice = string(text)
	new.Timestamp = data.Timestamp

	return new
}
