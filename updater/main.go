package main

import (
	"fmt"
	"os"
	"updater/cmd"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	switch os.Getenv("TYPE") {
	case "upload":
		cmd.Upload()
	case "download":
		cmd.Download()
	default:
		fmt.Println("Invalid type")
		return
	}
}
