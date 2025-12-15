package v2

import (
	"time"
)

func Get_status_update() Stats {
	status := Stats{
		Status:    "running",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return status
}
