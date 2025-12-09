package main

import (
	"fmt"
	"os"
	"strings"
	"updater/cmd"

	"github.com/joho/godotenv"
)

func isDev() bool {
	temp, _ := os.Executable()
	fmt.Println(temp)
	return strings.Contains(temp, os.TempDir()) || strings.Contains(temp, "\\Local\\go-build")
}

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
		cmd.Download(isDev())
	default:
		fmt.Println("Invalid type")
		return
	}
}
