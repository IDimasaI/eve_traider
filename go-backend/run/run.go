package run

import (
	"database/sql"
	"fmt"
	localfunc "local_server/local_func"
	updatepricelist "local_server/update_price_list"
	"net/http"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Run(isDev bool) error {
	//return nil
	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_TOKEN"))

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

	if isDev {
		//extractor := localfunc.NewExtractor()
		//client := localfunc.NewClient(http.DefaultClient)

		//file, _ := utils.ReadFile("E:\\dev\\eve\\eve_traider\\local_server\\items.json")
		//all_names, _, _ := extractor.Extract(file)

		//ids, err := client.Chunk_worker(all_names, 50)
		//if err != nil {
		//	fmt.Fprintf(os.Stderr, "failed to create ids: %s", err)
		//	os.Exit(1)
		//}
		//utils.WriteJson("E:\\dev\\eve\\eve_traider\\local_server\\ids.json", ids)

		// if err := localfunc.AddItemsFromJson(db, "E:\\dev\\eve\\eve_traider\\local_server\\ids.json"); err != nil {
		// 	fmt.Fprintf(os.Stderr, "failed to add items from json: %s", err)
		// 	os.Exit(1)
		// }

	}

	switch os.Getenv("VERSION_DB") {
	case "2":
		fmt.Println("Version 2")
		update_need(db, "false")
		updatePriceListV2 := updatepricelist.NewUpdatePriceListV2(db, http.DefaultClient)
		err := updatePriceListV2.UpdatePriceList()
		if err != nil {
			return fmt.Errorf("failed to update price list: %s", err)
		}
		update_need(db, "true")
	case "1":
		if err := createTable_if_not(db); err != nil {
			fmt.Fprintf(os.Stderr, "failed to query users: %s", err)
			os.Exit(1)
		}
		if need_update(db) {
			fmt.Println("Need update")
			err = update_need(db, "false")
			if err != nil {
				return fmt.Errorf("failed to update table: %s", err)
			}
			list := updatepricelist.NewUpdatePriceList(db)
			err = list.UpdatePriceList()
			if err != nil {
				return fmt.Errorf("failed to update price list: %s", err)
			}
			err = update_need(db, "true")
			if err != nil {
				return fmt.Errorf("failed to update table: %s", err)
			}
		} else {
			fmt.Println("No update needed")
		}
	default:
		fmt.Println("Unknown version")
	}
	return nil
}

func createTable_if_not(db *sql.DB) error {
	// Создаем таблицу если не существует
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS need_update (id INTEGER PRIMARY KEY, value TEXT, last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	// Проверяем, есть ли уже строка с id=1
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM need_update WHERE id = ?", 1).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking row existence: %w", err)
	}

	// Если строки нет - вставляем
	if count == 0 {
		create_all_tables_in_db(db)
	}

	return nil
}
func need_update(db *sql.DB) bool {

	var value string
	db.QueryRow("SELECT value FROM need_update WHERE id = 1").Scan(&value)
	if value == "true" {
		return true
	}
	return false
}

func update_need(db *sql.DB, value string) error {
	_, err := db.Exec("UPDATE need_update SET value = ?, last_updated = CURRENT_TIMESTAMP WHERE id = 1", value)
	if err != nil {
		return fmt.Errorf("error updating table: %w", err)
	}
	return nil
}

func create_all_tables_in_db(db *sql.DB) error {
	var err error
	//Тк я проверяю на создание и есть ли значения
	_, err = db.Exec("INSERT INTO need_update (id, value, last_updated) VALUES (?, ?, ?)", 1, "true", time.Now())
	if err != nil {
		return fmt.Errorf("error inserting row: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS all_ids (id INTEGER PRIMARY KEY, IdAndNames TEXT)")
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	jsonIdAndNames, err := localfunc.GetIdFromJson()
	if err != nil {
		return fmt.Errorf("error getting IdAndNames: %w", err)
	}
	_, err = db.Exec("INSERT INTO all_ids (id, IdAndNames) VALUES (?, ?)", 1, string(jsonIdAndNames))
	if err != nil {
		return fmt.Errorf("error inserting row: %w", err)
	}
	return nil
}
