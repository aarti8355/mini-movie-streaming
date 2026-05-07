package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedDefaultAdmin checks if an admin exists, and if not, creates one.
func SeedDefaultAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := DB.Exec("INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)",
			"Admin", "admin@minimovie.com", hashedPassword, "admin")
		if err != nil {
			log.Printf("Failed to seed default admin: %v", err)
		} else {
			log.Println("===============================================================")
			log.Println("Successfully seeded default admin user!")
			log.Println("Email: admin@minimovie.com")
			log.Println("Password: admin123")
			log.Println("===============================================================")
		}
	}
}
