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
	fmt.Println(temp)
	return strings.Contains(temp, os.TempDir()) || strings.Contains(temp, "\\Local\\go-build")
}

func main() {

	command := flag.String("command", "", "Command to execute")
	addr := flag.String("addr", "http://localhost:6969/api/v2/update_status", "Port to listen on")
	flag.Parse()

	if *command == "" {
		fmt.Println("Invalid type")
		fmt.Println("Usage: go-updater -command [upload|download]")
		return
	}

	switch *command {
	case "upload":
		cmd.Upload()
	case "download":
		cmd.Download(isDev(), *addr)
	default:
		fmt.Println("Invalid type")
		fmt.Println("Usage: go-updater -command [upload|download]")
		return
	}

}
