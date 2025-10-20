package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Импортируем pprof для профилирования
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cere6rum/MicroBlog2/internal/handlers"
	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
	"github.com/Cere6rum/MicroBlog2/internal/service"
)

func main() {
	// 1. Инициализация логгера
	appLogger, err := logger.NewLogger("app.log")
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer appLogger.Close()

	appLogger.Info("=== Запуск MicroBlog v1 ===")

	// 2. Создание очереди лайков (буфер 100, 3 воркера)
	likeQueue := queue.NewLikeQueue(100, 3)
	appLogger.Info("Очередь лайков создана (буфер: 100, воркеры: 3)")

	// 3. Создание сервиса бизнес-логики
	microBlogService := service.NewMicroBlogService(appLogger, likeQueue)
	appLogger.Info("Сервис MicroBlog инициализирован")

	// 4. Запуск обработчиков очереди лайков
	likeQueue.Start(microBlogService.ProcessLikeEvent)
	appLogger.Info("Воркеры очереди лайков запущены")

	// 5. Создание HTTP-обработчиков
	handler := handlers.NewMicroBlogHandler(microBlogService)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	appLogger.Info("HTTP-маршруты зарегистрированы")

	// 6. Запуск HTTP-сервера для профилирования на отдельном порту
	go func() {
		pprofAddr := ":6060"
		appLogger.Info(fmt.Sprintf("Профилирование pprof доступно на http://localhost%s/debug/pprof/", pprofAddr))
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			appLogger.Error(fmt.Sprintf("Ошибка запуска pprof сервера: %v", err))
		}
	}()

	// 7. Создание основного HTTP-сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. Запуск сервера в отдельной горутине
	go func() {
		appLogger.Info("HTTP-сервер запущен на :8080")
		fmt.Println("MicroBlog v1 запущен на http://localhost:8080")
		fmt.Println("Профилирование доступно на http://localhost:6060/debug/pprof/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error(fmt.Sprintf("Ошибка запуска сервера: %v", err))
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// 9. Graceful Shutdown - ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Получен сигнал завершения, начинаем graceful shutdown...")
	fmt.Println("Завершение работы...")

	// 10. Контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 11. Остановка HTTP-сервера
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error(fmt.Sprintf("Ошибка при остановке сервера: %v", err))
	} else {
		appLogger.Info("HTTP-сервер успешно остановлен")
	}

	// 12. Остановка очереди лайков
	likeQueue.Stop()
	appLogger.Info("Очередь лайков остановлена")

	appLogger.Info("=== MicroBlog v1 успешно завершен ===")
	fmt.Println("Приложение завершено")
}
