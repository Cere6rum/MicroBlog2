package models

// Post представляет пост в микроблоге
type Post struct {
	ID       int      `json:"id"`
	AuthorID int      `json:"author_id"`
	Author   string   `json:"author"`
	Content  string   `json:"content"`
	Likes    []string `json:"likes"` // Список пользователей, лайкнувших пост
}
