// scripts/build.go
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var cwd, _ = os.Getwd()

const NEW_VERSION_BUILD string = "v0.8"

func main() {

	create_webview_folders()
	create_resources()

	fmt.Println("Building EVE Price Tracker...")

	// Основной бекенд
	exec.Command("go", "build", "-C", "./go-backend", "-ldflags", "-H=windowsgui", "-o", "./../build/go-backend.exe").Run()
	exec.Command("go", "build", "-C", "./go-backend", "-ldflags", "", "-o", "./../build/go-backend_no_gui.exe").Run()
	// Обновлятор, с зашитой версией
	exec.Command("go", "build", "-C", "./go-updater", "-ldflags", "", "-o", "./../build/updater.exe").Run()

	build_launcher()

	config, err := os.Create("build/data/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer config.Close()
	_, _ = io.WriteString(config, `{"version": "`+NEW_VERSION_BUILD+`"}`)
	fmt.Println("Done!")
}

func create_webview_folders() {
	fmt.Println("Creating folders...")
	project_root := filepath.Join(cwd, "build")
	project_root = filepath.Join(project_root, "data")
	project_root = filepath.Join(project_root, "webview2")
	os.MkdirAll(project_root, os.ModePerm)
}

func create_resources() {
	fmt.Println("Creating resources...")
	project_root := filepath.Join(cwd, "build")
	project_root = filepath.Join(project_root, "data")
	project_root = filepath.Join(project_root, "resources")
	os.MkdirAll(project_root, os.ModePerm)
}

func build_launcher() {
	fmt.Println("🦀 Building Rust desktop launcher...")
	// Собираем проект
	cmd := exec.Command("cargo",
		"build",
		"--release",
		"--manifest-path",
		"rust-launcher/Cargo.toml")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Cargo build failed:", err)
	}

	// Копируем результат
	var exeName string
	if runtime.GOOS == "windows" {
		exeName = "launcher.exe"
	} else {
		exeName = "launcher"
	}

	src := filepath.Join("rust-launcher", "target", "release", exeName)
	dst := filepath.Join("build", exeName)

	if err := copyFile(src, dst); err != nil {
		log.Fatal("Copy failed:", err)
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
