package database

import (
	"damapp-server/config"
	"damapp-server/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fmt"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Sqlite_database), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Friendship{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
