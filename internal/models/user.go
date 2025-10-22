package models

// User представляет пользователя системы
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
