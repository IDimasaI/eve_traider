package localfunc

import (
	"encoding/json"
	"fmt"
	"local_server/utils"
)

type Items struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetIdFromJson() ([]byte, error) {
	file, err := utils.ReadJson[[]Items]("E:\\dev\\eve\\eve_traider\\local_server\\ids.json")
	if err != nil {
		return nil, fmt.Errorf("failed to parse json file: %s", err)
	}
	// var ids []string
	// var names []string

	// for _, item := range file {
	// 	ids = append(ids, fmt.Sprintf("%d", item.Id))
	// 	names = append(names, item.Name)
	// }

	jsonIdAndNames, err := json.Marshal(file)

	if err != nil {
		return nil, fmt.Errorf("failed \" to marshal json: %s", err)
	}

	//	utils.WriteFile("E:\\dev\\eve\\eve_traider\\local_server\\ids.txt", jsonBytes)

	return jsonIdAndNames, nil
}
