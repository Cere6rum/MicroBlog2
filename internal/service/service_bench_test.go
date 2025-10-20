package service

import (
	"fmt"
	"testing"

	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
)

// BenchmarkRegisterUser бенчмарк для регистрации пользователей
func BenchmarkRegisterUser(b *testing.B) {
	log, _ := logger.NewLogger("bench.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.RegisterUser(fmt.Sprintf("user%d", i))
	}
}

// BenchmarkCreatePost бенчмарк для создания постов
func BenchmarkCreatePost(b *testing.B) {
	log, _ := logger.NewLogger("bench.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	// Регистрируем одного пользователя
	service.RegisterUser("benchuser")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreatePost("benchuser", fmt.Sprintf("Пост номер %d", i))
	}
}

// BenchmarkGetAllPosts бенчмарк для получения всех постов
func BenchmarkGetAllPosts(b *testing.B) {
	log, _ := logger.NewLogger("bench.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(100, 2)
	service := NewMicroBlogService(log, likeQueue)

	// Создаем 100 постов
	service.RegisterUser("benchuser")
	for i := 0; i < 100; i++ {
		service.CreatePost("benchuser", fmt.Sprintf("Пост %d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAllPosts()
	}
}

// BenchmarkLikePost бенчмарк для лайков
func BenchmarkLikePost(b *testing.B) {
	log, _ := logger.NewLogger("bench.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(1000, 4)
	service := NewMicroBlogService(log, likeQueue)
	likeQueue.Start(service.ProcessLikeEvent)
	defer likeQueue.Stop()

	// Создаем пользователя и пост
	service.RegisterUser("author")
	service.RegisterUser("liker")
	service.CreatePost("author", "Бенчмарк пост")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.LikePost(1, "liker")
	}
}
