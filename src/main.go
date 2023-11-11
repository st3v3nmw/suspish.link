package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	limiter "github.com/ulule/limiter/v3"
	limiter_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiter_store "github.com/ulule/limiter/v3/drivers/store/memory"
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

	// Add middleware
	rate, _ := limiter.NewRateFromFormatted("60-S")
	store := limiter_store.NewStoreWithOptions(limiter.StoreOptions{Prefix: "limiter_gin", MaxRetry: 3})
	middleware := limiter_gin.NewMiddleware(limiter.New(store, rate))

	// Routers
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	router.Use(middleware)

	router.POST("/shorten", ShortenURL)
	router.GET("/*susURI", ResolveURL)

	router.Run()
}
