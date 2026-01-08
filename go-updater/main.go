package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"updater/cmd"
)

func isDev() bool {
	temp, _ := os.Executable()
	//fmt.Println(temp)
	return strings.Contains(temp, os.TempDir()) || strings.Contains(temp, "\\Local\\go-build")
}

func main() {

	command := flag.String("command", "", "Command to execute")
	repo := flag.String("repo", "", "User/repository") //"IDimasaI/eve_traider"
	addr := flag.String("addr", "", "http://*** Адрес для отправки пост статуса")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Println("Usage: go-updater -command [upload|download] -repo <repository_url> -addr <http://address>")
		fmt.Println("Example: go-updater -command download -repo IDimasaI/eve_traider -addr http://localhost:8080/api/update_status")
		fmt.Println("-addr опциональная команда для отправки post запросов о прогрессе.")
		return
	}

	if *command == "" {
		fmt.Println("Invalid type")
		fmt.Println("Usage: go-updater -command [upload|download]")
		return
	}
	if *repo == "" {
		fmt.Println("Invalid repository")
		fmt.Println("Usage: go-updater -command [upload|download] -repo <repository_url>")
		return
	}
	switch *command {
	case "upload":
		cmd.Upload()
	case "download":
		cmd.Download(isDev(), *addr, *repo)
	default:
		fmt.Println("Invalid type")
		fmt.Println("Usage: go-updater -command [upload|download]")
		return
	}

}
