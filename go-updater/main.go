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
		cmd.Download(isDev())
	default:
		fmt.Println("Invalid type")
		fmt.Println("Usage: go-updater -command [upload|download]")
		return
	}

}
