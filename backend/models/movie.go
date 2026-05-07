package models

type Movie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Genre         string  `json:"genre"`
	Description   string  `json:"description"`
	ReleaseYear   int     `json:"release_year"`
	ImageURL      string  `json:"image_url"`
	AverageRating float64 `json:"average_rating"`
	CreatedAt     string  `json:"created_at"`
}
