package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Cere6rum/MicroBlog2/internal/models"
)

// Logger - структура логгера с каналом
type Logger struct {
	logChan chan models.LogEvent
	done    chan struct{}
	file    *os.File
}

// NewLogger создает новый логгер
func NewLogger(logFile string) (*Logger, error) {
	// Открываем файл для логов (создаем если не существует)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл логов: %w", err)
	}

	logger := &Logger{
		logChan: make(chan models.LogEvent, 100), // Буферизованный канал
		done:    make(chan struct{}),
		file:    file,
	}

	// Запускаем горутину для обработки логов
	go logger.processLogs()

	return logger, nil
}

// processLogs обрабатывает события логирования из канала
func (l *Logger) processLogs() {
	for {
		select {
		case event := <-l.logChan:
			// Форматируем и записываем лог
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("[%s] [%s] %s\n", timestamp, event.Level, event.Message)

			// Пишем в файл
			if _, err := l.file.WriteString(logMessage); err != nil {
				log.Printf("Ошибка записи в лог: %v", err)
			}

			// Также выводим в stdout для удобства
			fmt.Print(logMessage)

		case <-l.done:
			// Завершаем обработку логов
			return
		}
	}
}

// Log отправляет событие в канал логирования
func (l *Logger) Log(level, message string) {
	l.logChan <- models.LogEvent{
		Level:   level,
		Message: message,
	}
}

// Info логирует информационное сообщение
func (l *Logger) Info(message string) {
	l.Log("INFO", message)
}

// Error логирует сообщение об ошибке
func (l *Logger) Error(message string) {
	l.Log("ERROR", message)
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(message string) {
	l.Log("DEBUG", message)
}

// Close закрывает логгер и освобождает ресурсы
func (l *Logger) Close() error {
	close(l.done)
	time.Sleep(100 * time.Millisecond) // Даем время обработать оставшиеся логи
	return l.file.Close()
}
