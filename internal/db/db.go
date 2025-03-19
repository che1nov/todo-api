package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-api/internal/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	log.Println("Successfully connected to SQLite")

	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatalf("Failed to perform auto-migration: %v", err)
	}
	log.Println("Database migration completed successfully")

	return db
}
