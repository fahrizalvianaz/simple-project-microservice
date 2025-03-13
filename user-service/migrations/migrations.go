package migrations

import (
	"bookstore-framework/internal/users"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&users.User{},
	)
	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
