package syncutils

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// SafeUserStorage - потокобезопасное хранилище пользователей
type SafeUserStorage struct {
	mu    sync.RWMutex
	users map[string]interface{} // интерфейс для гибкости, можно хранить *models.User
}

// NewSafeUserStorage создает новое потокобезопасное хранилище
func NewSafeUserStorage() *SafeUserStorage {
	return &SafeUserStorage{
		users: make(map[string]interface{}),
	}
}

// Set добавляет или обновляет пользователя
func (s *SafeUserStorage) Set(username string, user interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[username] = user
}

// Get возвращает пользователя по имени
func (s *SafeUserStorage) Get(username string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[username]
	return user, exists
}

// Exists проверяет существование пользователя
func (s *SafeUserStorage) Exists(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.users[username]
	return exists
}

// SafePostStorage - потокобезопасное хранилище постов
type SafePostStorage struct {
	mu    sync.RWMutex
	posts []interface{} // слайс постов
}

// NewSafePostStorage создает новое потокобезопасное хранилище постов
func NewSafePostStorage() *SafePostStorage {
	return &SafePostStorage{
		posts: make([]interface{}, 0),
	}
}

// Add добавляет новый пост
func (s *SafePostStorage) Add(post interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts = append(s.posts, post)
}

// GetAll возвращает все посты
func (s *SafePostStorage) GetAll() []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Возвращаем копию слайса для безопасности
	copied := make([]interface{}, len(s.posts))
	copy(copied, s.posts)
	return copied
}

// GetByIndex возвращает пост по индексу
func (s *SafePostStorage) GetByIndex(index int) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.posts) {
		return nil, false
	}
	return s.posts[index], true
}

// Len возвращает количество постов
func (s *SafePostStorage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.posts)
}

// AtomicCounter - атомарный счетчик для ID
type AtomicCounter struct {
	value int64
}

// NewAtomicCounter создает новый атомарный счетчик
func NewAtomicCounter(initial int64) *AtomicCounter {
	return &AtomicCounter{value: initial}
}

// Increment увеличивает счетчик и возвращает новое значение
func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Get возвращает текущее значение
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set устанавливает новое значение
func (c *AtomicCounter) Set(val int64) {
	atomic.StoreInt64(&c.value, val)
}

// SetByIndex заменяет пост по индексу, возвращает ошибку если индекс неверен
func (s *SafePostStorage) SetByIndex(index int, post interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.posts) {
		return fmt.Errorf("index out of range")
	}
	s.posts[index] = post
	return nil
}
