package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	dsn := "root:Kira@8355@tcp(db:3306)/mini_movie_db?parseTime=true"
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	DB = db
	log.Println("Database connection successfully established!")
	return nil
}
