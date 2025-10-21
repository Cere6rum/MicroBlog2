package service

import (
	"fmt"
	"testing"

	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
)

func BenchmarkRegisterUser(b *testing.B) {
	log, err := logger.NewLogger("bench.log")
	if err != nil {
		b.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()

	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := service.RegisterUser(fmt.Sprintf("user%d", i)); err != nil {
			b.Fatalf("Ошибка регистрации пользователя: %v", err)
		}
	}
}

func BenchmarkCreatePost(b *testing.B) {
	log, err := logger.NewLogger("bench.log")
	if err != nil {
		b.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()

	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	if _, err := service.RegisterUser("benchuser"); err != nil {
		b.Fatalf("Ошибка регистрации пользователя: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := service.CreatePost("benchuser", fmt.Sprintf("Пост номер %d", i)); err != nil {
			b.Fatalf("Ошибка создания поста: %v", err)
		}
	}
}

func BenchmarkGetAllPosts(b *testing.B) {
	log, err := logger.NewLogger("bench.log")
	if err != nil {
		b.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()

	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	if _, err := service.RegisterUser("benchuser"); err != nil {
		b.Fatalf("Ошибка регистрации пользователя: %v", err)
	}
	for i := 0; i < 100; i++ {
		if _, err := service.CreatePost("benchuser", fmt.Sprintf("Пост %d", i)); err != nil {
			b.Fatalf("Ошибка создания поста %d: %v", i, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAllPosts()
	}
}

func BenchmarkLikePost(b *testing.B) {
	log, err := logger.NewLogger("bench.log")
	if err != nil {
		b.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			b.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()

	likeQueue := queue.NewLikeQueue(1000, 4)
	service := NewMicroBlogService(log, likeQueue)
	likeQueue.Start(service.ProcessLikeEvent)
	defer likeQueue.Stop()

	if _, err := service.RegisterUser("author"); err != nil {
		b.Fatalf("Ошибка регистрации пользователя author: %v", err)
	}
	if _, err := service.RegisterUser("liker"); err != nil {
		b.Fatalf("Ошибка регистрации пользователя liker: %v", err)
	}
	if _, err := service.CreatePost("author", "Бенчмарк пост"); err != nil {
		b.Fatalf("Ошибка создания поста: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := service.LikePost(1, "liker"); err != nil {
			b.Fatalf("Ошибка лайка поста: %v", err)
		}
	}
}
