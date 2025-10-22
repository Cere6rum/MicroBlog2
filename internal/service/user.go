package service

import (
	"errors"
	"fmt"

	"github.com/Cere6rum/MicroBlog2/internal/models"
)

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

// GetUserByUsername возвращает пользователя по имени
func (s *MicroBlogService) GetUserByUsername(username string) (*models.User, error) {
	userInterface, exists := s.userStorage.Get(username)
	if !exists {
		return nil, errors.New("пользователь не найден")
	}
	return userInterface.(*models.User), nil
}
