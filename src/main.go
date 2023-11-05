package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Open the database connection
	db_dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(db_dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// Run migrations
	if err := DB.AutoMigrate(&Link{}); err != nil {
		panic(err.Error())
	}

	// Routers
	r := gin.Default()

	r.POST("/shorten", ShortenURL)
	r.GET("/*susURI", ResolveURL)

	r.Run()
}
