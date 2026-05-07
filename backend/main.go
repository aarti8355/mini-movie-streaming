package main

import (
	"backend/controller"
	"backend/db"
	"backend/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database Connection
	if err := db.InitDB(); err != nil {
		log.Fatalf("Could not initialize database connection: %v\n", err)
	}
	defer db.DB.Close()

	// Seed the default admin user
	db.SeedDefaultAdmin()

	// Initialize Gin app
	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // common frontend ports
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	api := r.Group("/api")
	{
		// Public Routes
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
		api.GET("/movies", controller.GetAllMovies)
		api.GET("/movies/:id", controller.GetMovieByID)
		api.GET("/movies/:id/reviews", controller.GetReviewsByMovie)

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Review routes
			protected.POST("/reviews", controller.AddReview)

			// Watchlist routes
			protected.GET("/watchlist", controller.GetUserWatchlist)
			protected.POST("/watchlist", controller.AddToWatchlist)
			protected.DELETE("/watchlist/:id", controller.RemoveFromWatchlist)

			// Admin Routes
			admin := protected.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.POST("/movies", controller.CreateMovie)
				admin.PUT("/movies/:id", controller.UpdateMovie)
				admin.DELETE("/movies/:id", controller.DeleteMovie)
				
				admin.GET("/reviews", controller.GetAllReviews)
				admin.DELETE("/reviews/:id", controller.DeleteReview)
			}
		}
	}

	log.Println("Starting API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
