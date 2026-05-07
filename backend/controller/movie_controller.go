package controller

import (
	"backend/db"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllMovies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	search := c.Query("search")
	genre := c.Query("genre")
	rating, _ := strconv.ParseFloat(c.Query("rating"), 64)
	sortParam := c.Query("sort")

	query := "SELECT id, title, genre, COALESCE(description, ''), release_year, COALESCE(image_url, ''), average_rating, created_at FROM movies WHERE 1=1"
	var args []interface{}

	if search != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if genre != "" {
		query += " AND genre LIKE ?"
		args = append(args, "%"+genre+"%")
	}
	if rating > 0 {
		query += " AND average_rating >= ?"
		args = append(args, rating)
	}

	if sortParam == "rating_desc" {
		query += " ORDER BY average_rating DESC"
	} else if sortParam == "rating_asc" {
		query += " ORDER BY average_rating ASC"
	} else {
		query += " ORDER BY id DESC"
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Description, &movie.ReleaseYear, &movie.ImageURL, &movie.AverageRating, &movie.CreatedAt)
		if err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	if movies == nil {
		movies = []models.Movie{}
	}

	// Also fetch total count for pagination
	countQuery := "SELECT COUNT(*) FROM movies WHERE 1=1"
	var countArgs []interface{}
	if search != "" {
		countQuery += " AND title LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}
	if genre != "" {
		countQuery += " AND genre LIKE ?"
		countArgs = append(countArgs, "%"+genre+"%")
	}
	if rating > 0 {
		countQuery += " AND average_rating >= ?"
		countArgs = append(countArgs, rating)
	}
	var total int
	db.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"page": page,
		"total_pages": totalPages,
		"total_movies": total,
	})
}

func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO movies (title, genre, description, release_year, image_url) VALUES (?, ?, ?, ?, ?)"
	res, err := db.DB.Exec(query, movie.Title, movie.Genre, movie.Description, movie.ReleaseYear, movie.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	id, _ := res.LastInsertId()
	movie.ID = int(id)
	
	c.JSON(http.StatusCreated, movie)
}

func GetMovieByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie models.Movie
	query := "SELECT id, title, genre, COALESCE(description, ''), release_year, COALESCE(image_url, ''), average_rating, created_at FROM movies WHERE id = ?"
	err = db.DB.QueryRow(query, id).Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Description, &movie.ReleaseYear, &movie.ImageURL, &movie.AverageRating, &movie.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func UpdateMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE movies SET title = ?, genre = ?, description = ?, release_year = ?, image_url = ? WHERE id = ?"
	_, err = db.DB.Exec(query, movie.Title, movie.Genre, movie.Description, movie.ReleaseYear, movie.ImageURL, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	movie.ID = id
	c.JSON(http.StatusOK, movie)
}

func DeleteMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	query := "DELETE FROM movies WHERE id = ?"
	res, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie successfully deleted"})
}
