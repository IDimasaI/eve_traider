package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"local_server/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Release struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`    //https://github.com/IDimasaI/eve_traider/releases/tag/v0
	Zipball_url string `json:"zipball_url"` //https://github.com/IDimasaI/eve_traider/archive/refs/tags/v0.zip
}

func getLatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	return &release, err
}

func downloadRelease(release *Release) error {
	fmt.Printf("Downloading release %s...\n", release.Zipball_url)
	resp, err := http.Get(release.Zipball_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("release-%s.zip", release.TagName), body, 0644)
	return err
}
func extractFile(zipFile *zip.File, targetPath string) error {
	fileInfo := zipFile.FileInfo()

	src, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Создаем файл с правильными правами
	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	// Сохраняем время модификации
	return os.Chtimes(targetPath, fileInfo.ModTime(), fileInfo.ModTime())
}

func unpackZip(zipPath, downloadPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	os.MkdirAll(downloadPath, 0755)

	foundBuildFolder := false
	buildPrefix := ""

	// Находим build папку
	for _, file := range reader.File {
		idx := strings.Index(file.Name, "/build/")
		if idx != -1 {
			buildPrefix = file.Name[:idx+len("/build/")]
			foundBuildFolder = true
			break
		}
	}

	if !foundBuildFolder {
		for _, file := range reader.File {
			if strings.HasPrefix(file.Name, "build/") {
				buildPrefix = "build/"
				foundBuildFolder = true
				break
			}
		}
	}

	if !foundBuildFolder {
		return fmt.Errorf("папка 'build' не найдена")
	}

	// Список файлов, которые могут быть запущены
	runningExecutables := map[string]bool{
		"go-backend_no_gui.exe": true,
		"go-backend.exe":        true,
		"launcher.exe":          true,
		"updater.exe":           true,
	}

	for _, file := range reader.File {
		if !strings.Contains(file.Name, buildPrefix) {
			continue
		}

		relPath := strings.TrimPrefix(file.Name, buildPrefix)
		if relPath == "" {
			continue
		}

		targetPath := filepath.Join(downloadPath, relPath)
		fileName := filepath.Base(targetPath)

		// Проверяем, является ли файл запущенным исполняемым файлом
		if runningExecutables[strings.ToLower(fileName)] {
			if os.Rename(targetPath, targetPath+".old") != nil {
				log.Printf("Ошибка переименования %s: %v", fileName, err)
			}
			
			
		}

		fileInfo := file.FileInfo()

		if fileInfo.IsDir() {
			os.MkdirAll(targetPath, fileInfo.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(targetPath), 0755)

		if err := extractFile(file, targetPath); err != nil {
			return fmt.Errorf("ошибка извлечения %s: %w", relPath, err)
		}

	}

	return nil
}
func deleteZipFile(release *Release) error {
	return os.RemoveAll(fmt.Sprintf("release-%s.zip", release.TagName))
}
func Download(isDev bool, addr string) {
	fmt.Println("Downloading...")

	http.PostForm(addr, url.Values{"update": {"start"}, "progress": {"start"}})

	var config_path string

	if isDev {
		config_path = "./../build/data/config.json"
	} else {
		config_path = "data/config.json"
	}

	// Создаем директорию
	os.MkdirAll(filepath.Dir(config_path), 0755)

	// Пробуем прочитать, если не получается - создаем новый
	config, err := utils.ReadJson[utils.Config](config_path)
	if err != nil {
		log.Printf("Error reading config: %v, creating default", err)
		config = utils.Config{Version: "0.0"}
		if err := utils.WriteJson[utils.Config](config_path, config); err != nil {
			log.Printf("Failed to create config: %v", err)
			return
		}
	}

	log.Println("Current version:", config.Version)

	release, err := getLatestRelease("IDimasaI", "eve_traider")
	if err != nil {
		http.PostForm(addr, url.Values{"update": {"error"}, "progress": {err.Error()}})
		fmt.Println("Error:", err)
		return
	}
	log.Println("Latest release:", release.Name, release.TagName)

	if release.TagName == config.Version {
		log.Println("Already up to date")
		http.PostForm(addr, url.Values{"update": {"finished"}, "progress": {"Already up to date"}})
		return
	}

	http.PostForm(addr, url.Values{"update": {"working"}, "progress": {"downloadRelease"}})
	err = downloadRelease(release)
	if err != nil {
		fmt.Println("Error:", err)
		http.PostForm(addr, url.Values{"update": {"error"}, "progress": {err.Error()}})
		return
	}

	var downloadPath string
	if isDev {
		downloadPath = "E:\\dev\\eve\\app_eve_traider\\go-updater"
	} else {
		downloadPath, _ = os.Executable()
		downloadPath = filepath.Join(filepath.Dir(downloadPath))
	}

	defer deleteZipFile(release)
	http.PostForm(addr, url.Values{"update": {"working"}, "progress": {"unpackRelease"}})
	err = unpackZip(fmt.Sprintf("release-%s.zip", release.TagName), downloadPath)
	if err != nil {
		fmt.Println("Error:", err)
		http.PostForm(addr, url.Values{"update": {"error"}, "progress": {err.Error()}})
		return
	}

	http.PostForm(addr, url.Values{"update": {"working"}, "progress": {"updateConfig"}})
	err = utils.WriteJson(config_path, utils.Config{Version: release.TagName})
	if err != nil {
		fmt.Println("Error:", err)
		http.PostForm(addr, url.Values{"update": {"error"}, "progress": {err.Error()}})
		return
	}

	http.PostForm(addr, url.Values{"update": {"finished"}, "progress": {"updateComplete"}})
	fmt.Printf("Latest release: %s (%s)\n", release.Name, release.TagName)
}
