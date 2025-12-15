package utils

import (
	"fmt"
	"strings"
	"time"
)

// Функция для парсинга и форматирования интервала из строки @every
func ParseEveryInterval(everyString string) (time.Duration, string, error) {
	// Проверяем, начинается ли строка с @every
	if !strings.HasPrefix(everyString, "@every ") {
		return 0, "", fmt.Errorf("неверный формат, ожидается @every <duration>")
	}

	// Извлекаем часть с продолжительностью
	durationStr := strings.TrimPrefix(everyString, "@every ")

	// Парсим продолжительность с помощью стандартной функции Go
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	// Форматируем продолжительность в удобочитаемый вид
	readable := formatDuration(duration)

	return duration, readable, nil
}

// Функция для преобразования time.Duration в удобочитаемую строку
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d час(ов)", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d минут(ы)", minutes))
	}
	if seconds > 0 && hours == 0 { // Секунды обычно показываем только если нет часов
		parts = append(parts, fmt.Sprintf("%d секунд(ы)", seconds))
	}

	return strings.Join(parts, " ")
}
