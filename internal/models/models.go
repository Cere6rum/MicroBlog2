package models

// LikeEvent представляет событие лайка для асинхронной обработки
type LikeEvent struct {
	PostID   int
	Username string
}

// LogEvent представляет событие для логирования
type LogEvent struct {
	Level   string // INFO, ERROR, DEBUG
	Message string
}
