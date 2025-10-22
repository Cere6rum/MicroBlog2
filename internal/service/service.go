package service

import (
	"github.com/Cere6rum/MicroBlog2/internal/logger"
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
