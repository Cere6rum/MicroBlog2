package service

import (
	"testing"

	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
)

// TestRegisterUser тестирует регистрацию пользователя
func TestRegisterUser(t *testing.T) {
	// Создаем тестовый логгер
	log, err := logger.NewLogger("test.log")
	if err != nil {
		t.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer log.Close()

	// Создаем очередь лайков
	likeQueue := queue.NewLikeQueue(10, 1)

	// Создаем сервис
	service := NewMicroBlogService(log, likeQueue)

	// Тест 1: успешная регистрация
	user, err := service.RegisterUser("testuser")
	if err != nil {
		t.Errorf("Ожидали успешную регистрацию, получили ошибку: %v", err)
	}
	if user == nil {
		t.Error("Ожидали пользователя, получили nil")
	}
	if user != nil && user.Username != "testuser" {
		t.Errorf("Ожидали имя 'testuser', получили '%s'", user.Username)
	}

	// Тест 2: повторная регистрация того же пользователя
	_, err = service.RegisterUser("testuser")
	if err == nil {
		t.Error("Ожидали ошибку при повторной регистрации")
	}

	// Тест 3: регистрация с пустым именем
	_, err = service.RegisterUser("")
	if err == nil {
		t.Error("Ожидали ошибку при регистрации с пустым именем")
	}
}

// TestCreatePost тестирует создание поста
func TestCreatePost(t *testing.T) {
	log, _ := logger.NewLogger("test.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Регистрируем пользователя
	service.RegisterUser("author")

	// Тест 1: создание поста
	post, err := service.CreatePost("author", "Мой первый пост")
	if err != nil {
		t.Errorf("Ошибка создания поста: %v", err)
	}
	if post == nil {
		t.Error("Ожидали пост, получили nil")
	}
	if post != nil && post.Content != "Мой первый пост" {
		t.Errorf("Ожидали содержимое 'Мой первый пост', получили '%s'", post.Content)
	}

	// Тест 2: создание поста несуществующим пользователем
	_, err = service.CreatePost("nonexistent", "Тест")
	if err == nil {
		t.Error("Ожидали ошибку при создании поста несуществующим пользователем")
	}

	// Тест 3: создание поста с пустым содержимым
	_, err = service.CreatePost("author", "")
	if err == nil {
		t.Error("Ожидали ошибку при создании поста с пустым содержимым")
	}
}

// TestGetAllPosts тестирует получение всех постов
func TestGetAllPosts(t *testing.T) {
	log, _ := logger.NewLogger("test.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Регистрируем пользователя и создаем посты
	service.RegisterUser("user1")
	service.CreatePost("user1", "Пост 1")
	service.CreatePost("user1", "Пост 2")

	// Получаем все посты
	posts := service.GetAllPosts()
	if len(posts) != 2 {
		t.Errorf("Ожидали 2 поста, получили %d", len(posts))
	}
}

// TestLikePost тестирует добавление лайка
func TestLikePost(t *testing.T) {
	log, _ := logger.NewLogger("test.log")
	defer log.Close()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Запускаем обработку очереди
	likeQueue.Start(service.ProcessLikeEvent)
	defer likeQueue.Stop()

	// Регистрируем пользователей и создаем пост
	service.RegisterUser("author")
	service.RegisterUser("liker")
	service.CreatePost("author", "Тестовый пост")

	// Тест 1: успешный лайк
	err := service.LikePost(1, "liker")
	if err != nil {
		t.Errorf("Ошибка при лайке поста: %v", err)
	}

	// Тест 2: лайк несуществующего поста
	err = service.LikePost(999, "liker")
	if err == nil {
		t.Error("Ожидали ошибку при лайке несуществующего поста")
	}

	// Тест 3: лайк несуществующим пользователем
	err = service.LikePost(1, "nonexistent")
	if err == nil {
		t.Error("Ожидали ошибку при лайке несуществующим пользователем")
	}
}
