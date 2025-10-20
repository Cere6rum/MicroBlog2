package service

import (
	"errors"
	"fmt"

	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/models"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
	"github.com/Cere6rum/MicroBlog2/internal/syncutils"
)

// MicroBlogService - основной сервис микроблога
type MicroBlogService struct {
	userStorage   *syncutils.SafeUserStorage
	postStorage   *syncutils.SafePostStorage
	userIDCounter *syncutils.AtomicCounter
	postIDCounter *syncutils.AtomicCounter
	likeQueue     *queue.LikeQueue
	logger        *logger.Logger
}

// NewMicroBlogService создает новый экземпляр сервиса
func NewMicroBlogService(log *logger.Logger, likeQueue *queue.LikeQueue) *MicroBlogService {
	return &MicroBlogService{
		userStorage:   syncutils.NewSafeUserStorage(),
		postStorage:   syncutils.NewSafePostStorage(),
		userIDCounter: syncutils.NewAtomicCounter(0),
		postIDCounter: syncutils.NewAtomicCounter(0),
		likeQueue:     likeQueue,
		logger:        log,
	}
}

// RegisterUser регистрирует нового пользователя
func (s *MicroBlogService) RegisterUser(username string) (*models.User, error) {
	if username == "" {
		s.logger.Error("Попытка регистрации с пустым именем пользователя")
		return nil, errors.New("имя пользователя не может быть пустым")
	}

	// Проверяем, существует ли пользователь
	if s.userStorage.Exists(username) {
		s.logger.Error(fmt.Sprintf("Пользователь %s уже существует", username))
		return nil, errors.New("пользователь уже существует")
	}

	// Создаем нового пользователя
	userID := int(s.userIDCounter.Increment())
	user := &models.User{
		ID:       userID,
		Username: username,
	}

	// Сохраняем в хранилище
	s.userStorage.Set(username, user)
	s.logger.Info(fmt.Sprintf("Зарегистрирован новый пользователь: %s (ID: %d)", username, userID))

	return user, nil
}

// CreatePost создает новый пост
func (s *MicroBlogService) CreatePost(username, content string) (*models.Post, error) {
	if content == "" {
		s.logger.Error("Попытка создания поста с пустым содержимым")
		return nil, errors.New("содержимое поста не может быть пустым")
	}

	// Проверяем существование пользователя
	userInterface, exists := s.userStorage.Get(username)
	if !exists {
		s.logger.Error(fmt.Sprintf("Пользователь %s не найден", username))
		return nil, errors.New("пользователь не найден")
	}

	user := userInterface.(*models.User)

	// Создаем новый пост
	postID := int(s.postIDCounter.Increment())
	post := &models.Post{
		ID:       postID,
		AuthorID: user.ID,
		Author:   user.Username,
		Content:  content,
		Likes:    make([]string, 0),
	}

	// Добавляем в хранилище
	s.postStorage.Add(post)
	s.logger.Info(fmt.Sprintf("Создан новый пост ID: %d от пользователя: %s", postID, username))

	return post, nil
}

// GetAllPosts возвращает все посты
func (s *MicroBlogService) GetAllPosts() []*models.Post {
	postsInterface := s.postStorage.GetAll()
	posts := make([]*models.Post, 0, len(postsInterface))

	for _, p := range postsInterface {
		if post, ok := p.(*models.Post); ok {
			posts = append(posts, post)
		}
	}

	s.logger.Debug(fmt.Sprintf("Запрошены все посты, количество: %d", len(posts)))
	return posts
}

// LikePost добавляет лайк к посту (асинхронно через очередь)
func (s *MicroBlogService) LikePost(postID int, username string) error {
	// Проверяем существование пользователя
	if !s.userStorage.Exists(username) {
		s.logger.Error(fmt.Sprintf("Пользователь %s не найден для лайка", username))
		return errors.New("пользователь не найден")
	}

	// Проверяем существование поста
	if postID <= 0 || postID > s.postStorage.Len() {
		s.logger.Error(fmt.Sprintf("Пост с ID %d не найден", postID))
		return errors.New("пост не найден")
	}

	// Отправляем событие в очередь для асинхронной обработки
	event := models.LikeEvent{
		PostID:   postID,
		Username: username,
	}
	s.likeQueue.Enqueue(event)
	s.logger.Info(fmt.Sprintf("Лайк от %s к посту %d добавлен в очередь", username, postID))

	return nil
}

// ProcessLikeEvent обрабатывает событие лайка (вызывается из очереди)
func (s *MicroBlogService) ProcessLikeEvent(event models.LikeEvent) error {
	// Получаем пост
	postInterface, exists := s.postStorage.GetByIndex(event.PostID - 1)
	if !exists {
		s.logger.Error(fmt.Sprintf("Пост с ID %d не найден при обработке лайка", event.PostID))
		return errors.New("пост не найден")
	}

	post := postInterface.(*models.Post)

	// Проверяем, не лайкал ли уже этот пользователь
	for _, liker := range post.Likes {
		if liker == event.Username {
			s.logger.Debug(fmt.Sprintf("Пользователь %s уже лайкнул пост %d", event.Username, event.PostID))
			return nil // Уже лайкнуто
		}
	}

	// Добавляем лайк
	post.Likes = append(post.Likes, event.Username)
	s.logger.Info(fmt.Sprintf("Лайк от %s к посту %d успешно обработан", event.Username, event.PostID))

	return nil
}

// GetUserByUsername возвращает пользователя по имени
func (s *MicroBlogService) GetUserByUsername(username string) (*models.User, error) {
	userInterface, exists := s.userStorage.Get(username)
	if !exists {
		return nil, errors.New("пользователь не найден")
	}
	return userInterface.(*models.User), nil
}
