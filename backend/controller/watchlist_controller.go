package controller

import (
	"backend/db"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToWatchlist(c *gin.Context) {
	var watchlist models.Watchlist
	if err := c.ShouldBindJSON(&watchlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	watchlist.UserID = userID.(int)

	query := "INSERT INTO watchlists (user_id, movie_id) VALUES (?, ?)"
	res, err := db.DB.Exec(query, watchlist.UserID, watchlist.MovieID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Movie is already in your watchlist"})
		return
	}

	id, _ := res.LastInsertId()
	watchlist.ID = int(id)

	c.JSON(http.StatusCreated, watchlist)
}

func GetUserWatchlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := `
		SELECT w.id, w.user_id, w.movie_id, w.created_at, m.title, m.image_url 
		FROM watchlists w
		JOIN movies m ON w.movie_id = m.id
		WHERE w.user_id = ?`
		
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch watchlist"})
		return
	}
	defer rows.Close()

	// Return a combined struct to include movie details
	type WatchlistItem struct {
		models.Watchlist
		MovieTitle    string `json:"movie_title"`
		MovieImageURL string `json:"movie_image_url"`
	}

	var items []WatchlistItem
	for rows.Next() {
		var item WatchlistItem
		err := rows.Scan(&item.ID, &item.UserID, &item.MovieID, &item.CreatedAt, &item.MovieTitle, &item.MovieImageURL)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	if items == nil {
		items = []WatchlistItem{}
	}

	c.JSON(http.StatusOK, items)
}

func RemoveFromWatchlist(c *gin.Context) {
	idStr := c.Param("id")
	watchlistID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid watchlist ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := "DELETE FROM watchlists WHERE id = ? AND user_id = ?"
	res, err := db.DB.Exec(query, watchlistID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from watchlist"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watchlist item not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist"})
}
