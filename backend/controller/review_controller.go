package controller

import (
	"backend/db"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if review.Comment != "" && len(review.Comment) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review comment must be at least 10 characters long"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	review.UserID = userID.(int)

	query := "INSERT INTO reviews (user_id, movie_id, rating, comment) VALUES (?, ?, ?, ?)"
	res, err := db.DB.Exec(query, review.UserID, review.MovieID, review.Rating, review.Comment)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already reviewed this movie"})
		return
	}

	id, _ := res.LastInsertId()
	review.ID = int(id)

	// Update average rating
	updateQuery := `
		UPDATE movies 
		SET average_rating = (SELECT AVG(rating) FROM reviews WHERE movie_id = ?) 
		WHERE id = ?`
	db.DB.Exec(updateQuery, review.MovieID, review.MovieID)

	c.JSON(http.StatusCreated, review)
}

func GetReviewsByMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	query := "SELECT id, user_id, movie_id, rating, comment, created_at FROM reviews WHERE movie_id = ?"
	rows, err := db.DB.Query(query, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.UserID, &review.MovieID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			continue
		}
		reviews = append(reviews, review)
	}

	if reviews == nil {
		reviews = []models.Review{}
	}

	c.JSON(http.StatusOK, reviews)
}

func GetAllReviews(c *gin.Context) {
	query := `
		SELECT r.id, r.user_id, r.movie_id, r.rating, r.comment, r.created_at, u.name as user_name, m.title as movie_title
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN movies m ON r.movie_id = m.id
		ORDER BY r.id DESC`
	
	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all reviews"})
		return
	}
	defer rows.Close()

	type AdminReview struct {
		models.Review
		UserName   string `json:"user_name"`
		MovieTitle string `json:"movie_title"`
	}

	var reviews []AdminReview
	for rows.Next() {
		var r AdminReview
		if err := rows.Scan(&r.ID, &r.UserID, &r.MovieID, &r.Rating, &r.Comment, &r.CreatedAt, &r.UserName, &r.MovieTitle); err == nil {
			reviews = append(reviews, r)
		}
	}

	if reviews == nil {
		reviews = []AdminReview{}
	}

	c.JSON(http.StatusOK, reviews)
}

func DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	query := "DELETE FROM reviews WHERE id = ?"
	res, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
