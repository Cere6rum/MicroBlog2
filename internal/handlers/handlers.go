package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Cere6rum/MicroBlog2/internal/service"
)

// MicroBlogHandler - обработчик HTTP-запросов
type MicroBlogHandler struct {
	service *service.MicroBlogService
}

// NewMicroBlogHandler создает новый обработчик
func NewMicroBlogHandler(svc *service.MicroBlogService) *MicroBlogHandler {
	return &MicroBlogHandler{
		service: svc,
	}
}

// RegisterRoutes регистрирует все маршруты
func (h *MicroBlogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", h.RegisterUser)
	mux.HandleFunc("/posts", h.PostsHandler)
	mux.HandleFunc("/posts/", h.LikePostHandler) // /posts/{id}/like
}

// RegisterUser обрабатывает POST /register
func (h *MicroBlogHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Регистрируем пользователя
	user, err := h.service.RegisterUser(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

// PostsHandler обрабатывает GET /posts и POST /posts
func (h *MicroBlogHandler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllPosts(w, r)
	case http.MethodPost:
		h.CreatePost(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// GetAllPosts обрабатывает GET /posts
func (h *MicroBlogHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts := h.service.GetAllPosts()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

// CreatePost обрабатывает POST /posts
func (h *MicroBlogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	post, err := h.service.CreatePost(req.Username, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}

// LikePostHandler обрабатывает POST /posts/{id}/like
func (h *MicroBlogHandler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим ID из URL (например, /posts/1/like)
	// Простой парсинг без использования gorilla/mux
	var postID int
	_, err := fmt.Sscanf(r.URL.Path, "/posts/%d/like", &postID)
	if err != nil {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	// Парсим JSON с именем пользователя
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Добавляем лайк
	if err := h.service.LikePost(postID, req.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"message":"Лайк успешно добавлен"}`)); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}
