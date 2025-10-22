package service

import (
	"errors"
	"fmt"

	"github.com/Cere6rum/MicroBlog2/internal/models"
)

// CreatePost создает новый пост
func (s *MicroBlogService) CreatePost(username, content string) (*models.Post, error) {
	if content == "" {
		s.logger.Error("Попытка создания поста с пустым содержимым")
		return nil, errors.New("содержимое поста не может быть пустым")
	}

	// Проверяем существование пользователя
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Пользователь %s не найден: %v", username, err))
		return nil, errors.New("пользователь не найден")
	}

	// Создаем новый пост
	postID := int(s.postIDCounter.Increment())
	post := &models.Post{
		ID:       postID,
		AuthorID: user.ID,
		Author:   user.Username,
		Content:  content,
		Likes:    make([]string, 0),
	}

	// Добавляем в репозиторий
	if err := s.postRepo.Create(post); err != nil {
		s.logger.Error(fmt.Sprintf("Ошибка при создании поста: %v", err))
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("Создан новый пост ID: %d от пользователя: %s", postID, username))

	return post, nil
}

// GetAllPosts возвращает все посты
func (s *MicroBlogService) GetAllPosts() ([]*models.Post, error) {
	posts := s.postRepo.List()
	s.logger.Debug(fmt.Sprintf("Запрошены все посты, количество: %d", len(posts)))
	return posts, nil
}

// LikePost добавляет лайк к посту (асинхронно через очередь)
func (s *MicroBlogService) LikePost(postID int, username string) error {
	// Проверяем существование пользователя
	if !s.userRepo.Exists(username) {
		s.logger.Error(fmt.Sprintf("Пользователь %s не найден для лайка", username))
		return errors.New("пользователь не найден")
	}

	// Проверяем существование поста
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Пост с ID %d не найден: %v", postID, err))
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
	post, err := s.postRepo.GetByID(event.PostID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Пост с ID %d не найден при обработке лайка: %v", event.PostID, err))
		return errors.New("пост не найден")
	}

	// Проверяем, не лайкал ли уже этот пользователь
	for _, liker := range post.Likes {
		if liker == event.Username {
			s.logger.Debug(fmt.Sprintf("Пользователь %s уже лайкнул пост %d", event.Username, event.PostID))
			return nil // Уже лайкнуто
		}
	}

	// Добавляем лайк и обновляем в репозитории
	post.Likes = append(post.Likes, event.Username)
	if err := s.postRepo.Update(post); err != nil {
		s.logger.Error(fmt.Sprintf("Ошибка при обновлении поста после лайка: %v", err))
		return err
	}

	s.logger.Info(fmt.Sprintf("Лайк от %s к посту %d успешно обработан", event.Username, event.PostID))
	return nil
}
