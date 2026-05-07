package models

type Watchlist struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	MovieID   int    `json:"movie_id"`
	CreatedAt string `json:"created_at"`
}
