package models

// SearchEvent поисковое событие, которое приходит от сервиса поиска
type SearchEvent struct {
	Query     string `json:"query"`     // поисковый запрос
	Timestamp int64  `json:"timestamp"` // время события (Unix timestamp)
	UserID    string `json:"user_id"`   // идентификатор пользователя
	IP        string `json:"ip"`        // IP-адрес для защиты от накруток
}
