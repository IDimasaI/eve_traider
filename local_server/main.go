package main

import (
	"context"
	"fmt"
	"local_server/run"
	"local_server/web"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Создаем контекст для graceful shutdown
	ctx_web, cancel_web := context.WithCancel(context.Background())
	defer cancel_web()

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup
	errorChan := make(chan error, 1)

	// Запускаем HTTP сервер в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()

		web.StartHTTPServer(ctx_web, errorChan)
	}()
	// Режим отладки
	if len(os.Args) > 1 && os.Args[1] == "-debug" {
		fmt.Println("Режим отладки активирован")

		// Ожидаем сигнал завершения в режиме отладки
		// if err := run.Run(true); err != nil {
		// 	fmt.Println("Ошибка при запуске режима отладки:", err)
		// 	fmt.Println("Режим отладки завершен")
		// }
		fmt.Println("Режим отладки RUN завершен")
		waitForShutdown()
		return
	}

	// Основной режим работы
	fmt.Println("EVE Online магазин парсер запущен!")

	if os.Getenv("CRON_JOB") != "false" {
		fmt.Println("> Расписание: " + os.Getenv("CRON_JOB"))

		// Создаем и настраиваем cron
		c := cron.New()
		var entryID cron.EntryID

		// Добавляем основную задачу по расписанию
		entryID, err = c.AddFunc(os.Getenv("CRON_JOB"), func() {
			go executeWithRetry("по расписанию", entryID, c, false)
		})
		if err != nil {
			log.Fatalf("❌ Ошибка настройки CRON: %v", err)
		}

		// Запускаем планировщик
		c.Start()
		defer c.Stop() // Останавливаем при завершении

		// Показываем информацию о следующем запуске
		entry := c.Entry(entryID)
		nextRunTime := entry.Schedule.Next(time.Now())
		timeUntilNext := time.Until(nextRunTime)

		fmt.Printf("> Следующий запрос через: %v\n> В %s\n", timeUntilNext.Round(time.Second), nextRunTime.Format("15:04:05 2006-01-02"))
	}
	fmt.Printf("💡 Приложение работает в фоне\n\n")

	go func() {
		select {
		case err := <-errorChan:
			log.Printf("❌ Критическая ошибка: %v", err)
			cancel_web() // Инициируем graceful shutdown
			fmt.Println("Все остальное вроде работает")
		case <-ctx_web.Done():

		}
	}()
	// Ожидаем сигнал завершения
	waitForShutdown()
	cancel_web()

	// Ждем завершения всех горутин
	wg.Wait()
	fmt.Println("👋 Приложение завершено")
}

// startHTTPServer запускает HTTP сервер

// waitForShutdown ожидает сигнал завершения
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n🛑 Получен сигнал завершения...")

}

// executeWithRetry выполняет парсинг с повторными попытками при ошибках
func executeWithRetry(reason string, entryID cron.EntryID, c *cron.Cron, isDev bool) {
	fmt.Printf("\n🔍 Запускаем парсинг (%s)...\n", reason)
	time_start := time.Now()

	for attempt := 1; attempt <= 3; attempt++ {
		fmt.Printf("🔄 Попытка %d...\n", attempt)

		if err := run.Run(isDev); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Попытка %d failed: %s\n", attempt, err)

			if attempt < 3 {
				fmt.Printf("⏳ Ждем 5 секунд перед повторной попыткой...\n")
				time.Sleep(5 * time.Second)
			}
		} else {
			// Успешное выполнение
			duration := time.Since(time_start)
			fmt.Println("★ Парсинг завершен успешно!")
			fmt.Printf("⧖ Время выполнения: %v\n", duration)

			// Показываем информацию о следующем запуске, если доступно
			if entryID != 0 && c != nil {
				entry := c.Entry(entryID)
				nextRunTime := entry.Schedule.Next(time.Now())
				timeUntilNext := time.Until(nextRunTime)

				fmt.Printf("> Следующий запрос через: %v\n", timeUntilNext.Round(time.Second))
				fmt.Printf("> В %s\n", nextRunTime.Format("15:04:05 2006-01-02"))
			}

			return
		}
	}

	// Все попытки провалились
	fmt.Fprintf(os.Stderr, "💥 Все 3 попытки парсинга провалились\n")

	// Даже при ошибке показываем следующее время запуска
	if entryID != 0 && c != nil {
		entry := c.Entry(entryID)
		nextRunTime := entry.Schedule.Next(time.Now())
		timeUntilNext := time.Until(nextRunTime)

		fmt.Printf("🔄 Следующая автоматическая попытка через: %v\n", timeUntilNext.Round(time.Second))
	}
}
