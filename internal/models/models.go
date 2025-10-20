package models

// User представляет пользователя системы
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Post представляет пост в микроблоге
type Post struct {
	ID       int      `json:"id"`
	AuthorID int      `json:"author_id"`
	Author   string   `json:"author"`
	Content  string   `json:"content"`
	Likes    []string `json:"likes"` // Список пользователей, лайкнувших пост
}

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
