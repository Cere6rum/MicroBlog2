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
	defer func() {
		if err := log.Close(); err != nil {
			t.Errorf("ошибка закрытия логгера: %v", err)
		}
	}()

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
	defer func() {
		if err := log.Close(); err != nil {
			t.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Регистрируем пользователя
	user, err := service.RegisterUser("author")
	if err != nil {
		t.Fatalf("Ошибка регистрации пользователя: %v", err)
	}
	if user == nil {
		t.Fatal("Ожидали пользователя после регистрации, получили nil")
	}

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
	defer func() {
		if err := log.Close(); err != nil {
			t.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Регистрируем пользователя и создаем посты
	user, err := service.RegisterUser("user1")
	if err != nil {
		t.Fatalf("Ошибка регистрации пользователя: %v", err)
	}
	if user == nil {
		t.Fatalf("Ожидали пользователя после регистрации, получили nil")
	}

	post1, err := service.CreatePost("user1", "Пост 1")
	if err != nil {
		t.Fatalf("Ошибка создания первого поста: %v", err)
	}
	if post1 == nil {
		t.Fatalf("Ожидали первый пост, получили nil")
	}

	post2, err := service.CreatePost("user1", "Пост 2")
	if err != nil {
		t.Fatalf("Ошибка создания второго поста: %v", err)
	}
	if post2 == nil {
		t.Fatalf("Ожидали второй пост, получили nil")
	}

	// Получаем все посты
	posts, err := service.GetAllPosts()
	if err != nil {
		t.Fatalf("Ожидали успешное получение постов, получили ошибку: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("Ожидали 2 поста, получили %d", len(posts))
	}
}

// TestLikePost тестирует добавление лайка
func TestLikePost(t *testing.T) {
	log, _ := logger.NewLogger("test.log")
	defer func() {
		if err := log.Close(); err != nil {
			t.Errorf("Ошибка закрытия логгера: %v", err)
		}
	}()
	likeQueue := queue.NewLikeQueue(10, 1)
	service := NewMicroBlogService(log, likeQueue)

	// Запускаем обработку очереди
	likeQueue.Start(service.ProcessLikeEvent)
	defer likeQueue.Stop()

	// Регистрируем пользователей и создаем пост
	user1, err := service.RegisterUser("author")
	if err != nil {
		t.Fatalf("Ошибка регистрации пользователя author: %v", err)
	}
	if user1 == nil {
		t.Fatal("Ожидали пользователя author после регистрации, получили nil")
	}

	user2, err := service.RegisterUser("liker")
	if err != nil {
		t.Fatalf("Ошибка регистрации пользователя liker: %v", err)
	}
	if user2 == nil {
		t.Fatal("Ожидали пользователя liker после регистрации, получили nil")
	}

	post, err := service.CreatePost("author", "Тестовый пост")
	if err != nil {
		t.Fatalf("Ошибка создания поста: %v", err)
	}
	if post == nil {
		t.Fatal("Ожидали пост после создания, получили nil")
	}
	postID := post.ID

	// Тест 1: успешный лайк
	err = service.LikePost(postID, "liker")
	if err != nil {
		t.Errorf("Ошибка при лайке поста: %v", err)
	}

	// Тест 2: лайк несуществующего поста
	err = service.LikePost(999, "liker")
	if err == nil {
		t.Error("Ожидали ошибку при лайке несуществующего поста")
	}

	// Тест 3: лайк несуществующим пользователем
	err = service.LikePost(postID, "nonexistent")
	if err == nil {
		t.Error("Ожидали ошибку при лайке несуществующим пользователем")
	}
}
