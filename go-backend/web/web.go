package web

import (
	"compress/gzip"
	"context"
	"crypto/md5"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"local_server/web/api"
	v2 "local_server/web/api/v2"
	proxy "local_server/web/proxys"
	specapi "local_server/web/spec_api"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func StartHTTPServer(ctx context.Context, errorChan chan error) {
	mux := http.NewServeMux()

	cache := proxy.NewProxy(proxy.InMemory)

	url := fmt.Sprintf("%s?authToken=%s", os.Getenv("TURSO_URL"), os.Getenv("TURSO_TOKEN"))
	db, err := sql.Open("libsql", url)
	if err != nil {
		errorChan <- err
		return
	}
	defer db.Close()

	api2 := v2.NewApi2(db)
	// Основной маршрут
	// любой кроме api/* и static/*
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/static/") {
			html, err := os.ReadFile("./web/pages/index.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write(html)
		}
	})

	//pages
	{
		staticFileServer := http.FileServer(http.Dir("./web/static"))
		mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

		// Маршрут для здоровья
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(cache.GetAll())

		})
	}

	//endpoints
	{
		mux.HandleFunc("/api/all_prices", func(w http.ResponseWriter, r *http.Request) {
			// Получаем данные
			prices, err := api.All_prices()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			// Сериализуем в JSON
			jsonData, err := json.Marshal(prices)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			// Вычисляем ETag
			hash := md5.Sum(jsonData)
			etag := fmt.Sprintf("\"%x\"", hash)

			// Проверяем If-None-Match
			if match := r.Header.Get("If-None-Match"); match != "" {
				// Убираем "W/" если есть (для слабых ETag)
				cleanMatch := strings.TrimPrefix(match, "W/")
				if cleanMatch == etag {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}

			// Проверяем, поддерживает ли клиент gzip
			acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

			// Устанавливаем основные заголовки
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
			w.Header().Set("ETag", etag)
			w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

			// Если клиент поддерживает gzip, сжимаем ответ
			if acceptsGzip {
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")

				w.WriteHeader(http.StatusOK)
				gz := gzip.NewWriter(w)
				defer gz.Close()

				if _, err := gz.Write(jsonData); err != nil {
					// В случае ошибки сжатия, логируем, но не прерываем выполнение
					fmt.Printf("Gzip error: %v\n", err)
				}
			} else {
				// Без сжатия
				w.WriteHeader(http.StatusOK)
				w.Write(jsonData)
			}
		})

		mux.HandleFunc("/api/all_items", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
			w.Header().Set("Expires", time.Now().Add(time.Hour).Format(time.RFC1123))
			items, err := api.All_items()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(items)
		})

		mux.HandleFunc("/api/v2/get_all_prices", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			var items []v2.PriceEntry
			var err error = nil

			if data, _ := proxy.GetAs[[]v2.PriceEntry](cache, "get_all_prices"); data != nil {
				items = data
				//	fmt.Println("Cache hit")
			} else {
				items, err = api2.Get_all_prices()
				//fmt.Println("Cache miss")
				go cache.Set("get_all_prices", items, 1*time.Hour)
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(items)
		})
		mux.HandleFunc("/api/v2/get_all_items", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			var items []v2.Item
			var err error = nil

			if data, _ := proxy.GetAs[[]v2.Item](cache, "get_all_items"); data != nil {
				items = data
				//	fmt.Println("Cache hit")
			} else {
				items, err = api2.Get_all_items()
				//	fmt.Println("Cache miss")
				go cache.Set("get_all_items", items, 1*time.Hour)
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(items)
		})
		mux.HandleFunc("/api/v2/get_prices", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			id := r.URL.Query().Get("id")

			items, err := api2.Get_prices(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(items)
		})

	}

	//техническая информация
	{
		Update_Status := v2.New_Update_Status()
		mux.HandleFunc("/api/v2/update_status", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				if err := Update_Status.Update_Status(w, r); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
					return
				}
			} else {
				w.Header().Set("Content-Type", "application/json")

				status := Update_Status.Get_Update_Status()

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(status)
			}
		})
	}

	//Специфичные для платформы эндпоинты
	{
		spe := specapi.CreatePlatform(os.Getenv("PLATFORM") == "cloud")
		spe.Cloud.Start()
	}

	server := &http.Server{
		Addr:    strings.ReplaceAll(os.Getenv("WEB_SERVER_ADDR"), "http://", ""),
		Handler: mux,
	}

	// Запуск сервера в отдельной горутине

	fmt.Println("🌐 HTTP сервер запущен на порту ", os.Getenv("WEB_SERVER_ADDR"), "\n ")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorChan <- fmt.Errorf("HTTP сервер упал: %v", err)
		}
	}()
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Ошибка при остановке HTTP сервера: %v\n", err)
	} else {
		fmt.Println("🌐 HTTP сервер остановлен")
	}
}
