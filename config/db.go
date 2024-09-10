package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func DB() *gorm.DB {
	dbOnce.Do(func() {
		// Load environment variables
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")

		// Construct the DSN
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=internal port=5432 sslmode=disable", host, user, password)
		var err error
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}

	})
	return dbInstance
}
