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
	if s.userRepo.Exists(username) {
		s.logger.Error(fmt.Sprintf("Пользователь %s уже существует", username))
		return nil, errors.New("пользователь уже существует")
	}

	// Создаем нового пользователя
	userID := int(s.userIDCounter.Increment())
	user := &models.User{
		ID:       userID,
		Username: username,
	}

	// Сохраняем в репозитории
	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error(fmt.Sprintf("Ошибка при создании пользователя: %v", err))
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("Зарегистрирован новый пользователь: %s (ID: %d)", username, userID))

	return user, nil
}

// GetUserByUsername возвращает пользователя по имени
func (s *MicroBlogService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}
