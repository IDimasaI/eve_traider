package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	localfunc "local_server/local_func"
	"os"
)

func All_items() ([]localfunc.Items, error) {
	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_TOKEN"))

	db, err := sql.Open("libsql", url)
	if err != nil {
		return []localfunc.Items{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT IdAndNames FROM all_ids WHERE id = 1")
	if err != nil {
		return []localfunc.Items{}, err
	}
	defer rows.Close()

	var IdAndNames []localfunc.Items
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		if err != nil {
			return []localfunc.Items{}, err
		}
		err = json.Unmarshal([]byte(row), &IdAndNames)
		if err != nil {
			return []localfunc.Items{}, err
		}
	}

	return IdAndNames, nil
}
