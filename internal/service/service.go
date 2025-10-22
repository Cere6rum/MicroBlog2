package service

import (
	"github.com/Cere6rum/MicroBlog2/internal/logger"
	"github.com/Cere6rum/MicroBlog2/internal/queue"
	"github.com/Cere6rum/MicroBlog2/internal/repository"
	"github.com/Cere6rum/MicroBlog2/internal/syncutils"
)

// MicroBlogService - основной сервис микроблога
type MicroBlogService struct {
	userRepo      repository.UserRepository
	postRepo      repository.PostRepository
	userIDCounter *syncutils.AtomicCounter
	postIDCounter *syncutils.AtomicCounter
	likeQueue     *queue.LikeQueue
	logger        *logger.Logger
}

// NewMicroBlogService создает новый экземпляр сервиса (обратная совместимость)
func NewMicroBlogService(log *logger.Logger, likeQueue *queue.LikeQueue) *MicroBlogService {
	ur := repository.NewInMemoryUserRepo()
	pr := repository.NewInMemoryPostRepo()
	return NewMicroBlogServiceWithRepos(log, likeQueue, ur, pr)
}

// NewMicroBlogServiceWithRepos создаёт сервис с подставными репозиториями (удобно для тестов)
func NewMicroBlogServiceWithRepos(log *logger.Logger, likeQueue *queue.LikeQueue, ur repository.UserRepository, pr repository.PostRepository) *MicroBlogService {
	return &MicroBlogService{
		userRepo:      ur,
		postRepo:      pr,
		userIDCounter: syncutils.NewAtomicCounter(0),
		postIDCounter: syncutils.NewAtomicCounter(0),
		likeQueue:     likeQueue,
		logger:        log,
	}
}
