package models

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // Don't expose password in JSON
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
}
