package v2

import (
	"os"
	"time"
)

func Get_status_update() Stats {
	status := Stats{
		Status:        "running",
		Timestamp:     time.Now().Format(time.RFC3339),
		CURRENT_BUILD: os.Getenv("CURRENT_BUILD"),
	}

	return status
}
