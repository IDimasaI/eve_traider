package v2

import (
	"fmt"
	"local_server/utils"
	"net/http"
	"time"
)

type Status struct {
	Current_Version string
	Status          string
	Progress        string
	Timestamp       int64
}

func New_Update_Status() Status {
	var config_path string
	var config utils.Config
	var err error
	if utils.IsDev() {
		config_path = "./../build/data/config.json"
		config, err = utils.ReadJson[utils.Config](config_path)
		if err != nil {
			fmt.Println("Error:", err)
			return Status{}
		}
	} else {
		config_path = "data/config.json"
		config, err = utils.ReadJson[utils.Config](config_path)
		if err != nil {
			fmt.Println("Error:", err)
			return Status{}
		}
	}

	return Status{
		Current_Version: config.Version,
		Status:          "",
		Progress:        "",
		Timestamp:       time.Now().Unix(),
	}
}

func (s *Status) Get_Update_Status() Status {

	status := Status{
		Current_Version: s.Current_Version,
		Status:          s.Status,
		Progress:        s.Progress,
		Timestamp:       time.Now().Unix(),
	}

	return status
}

func (s *Status) Update_Status(w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	form := r.PostForm

	s.Status = form.Get("update")
	s.Progress = form.Get("progress")
	return nil
}
