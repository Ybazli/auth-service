package models

type Session struct {
	UserID       uint   `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IP           string `json:"ip"`
	CreatedAt    string `json:"created_at"`
}
