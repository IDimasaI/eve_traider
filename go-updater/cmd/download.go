package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"local_server/utils"
	"log"
	"net/http"
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
type Config struct {
	Version string `json:"version"`
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

func unpackZip(zipPath, downloadPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть архив: %w", err)
	}
	defer reader.Close()

	// Создаем целевую папку, если её нет
	if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
		return fmt.Errorf("не удалось создать папку назначения: %w", err)
	}

	foundBuildFolder := false
	buildPrefix := ""

	// Сначала находим правильный префикс пути к папке build
	for _, file := range reader.File {
		// Ищем файлы/папки, содержащие "/build/" в пути
		idx := strings.Index(file.Name, "/build/")
		if idx != -1 {
			buildPrefix = file.Name[:idx+len("/build/")]
			foundBuildFolder = true
			break
		}
	}

	// Если не нашли "/build/", возможно папка build в корне архива?
	if !foundBuildFolder {
		// Проверяем, есть ли файлы, начинающиеся с "build/"
		for _, file := range reader.File {
			if strings.HasPrefix(file.Name, "build/") {
				buildPrefix = "build/"
				foundBuildFolder = true
				break
			}
		}
	}

	if !foundBuildFolder {
		return fmt.Errorf("папка 'build' не найдена в архиве")
	}

	// Теперь извлекаем только содержимое папки build
	for _, file := range reader.File {
		// Пропускаем файлы вне папки build
		if !strings.Contains(file.Name, buildPrefix) || strings.Contains(file.Name, "updater") {
			continue
		}

		// Получаем относительный путь внутри build
		relPath := strings.TrimPrefix(file.Name, buildPrefix)

		// Пропускаем саму папку build (пустой относительный путь)
		if relPath == "" {
			continue
		}

		targetPath := filepath.Join(downloadPath, relPath)

		// Обрабатываем папки
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("не удалось создать папку: %w", err)
			}
			continue
		}

		// Обрабатываем файлы
		// Создаем родительские папки для файла
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return fmt.Errorf("не удалось создать родительские папки: %w", err)
		}

		// Открываем файл в архиве
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("не удалось открыть файл в архиве: %w", err)
		}

		// Создаем файл на диске
		dstFile, err := os.Create(targetPath)
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("не удалось создать файл: %w", err)
		}

		// Копируем содержимое
		_, err = io.Copy(dstFile, srcFile)

		// Закрываем файлы в правильном порядке
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			return fmt.Errorf("ошибка копирования: %w", err)
		}
	}

	return nil
}
func deleteZipFile(release *Release) error {
	return os.RemoveAll(fmt.Sprintf("release-%s.zip", release.TagName))
}
func Download(isDev bool) {
	fmt.Println("Downloading...")

	var config_path string

	if isDev {
		config_path = "./../build/data/config.json"
	} else {
		config_path = "./data/config.json"
	}

	var config Config
	config, err := utils.ReadJson[Config](config_path)
	if err != nil {
		if err == os.ErrNotExist {
			log.Println("Config file not found")
			config = Config{Version: "0.0"}
			os.MkdirAll(filepath.Dir(config_path), 0755)
			utils.WriteJson[Config](config_path, config)
		} else {
			log.Println("Error reading config file:", err)
			return
		}
	}
	log.Println("Current version:", config.Version)

	release, err := getLatestRelease("IDimasaI", "eve_traider")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Println("Latest release:", release.Name, release.TagName)

	if release.TagName == config.Version {
		log.Println("Already up to date")
		return
	}

	err = downloadRelease(release)
	if err != nil {
		fmt.Println("Error:", err)
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
	err = unpackZip(fmt.Sprintf("release-%s.zip", release.TagName), downloadPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Latest release: %s (%s)\n", release.Name, release.TagName)
}
