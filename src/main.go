package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	limiter "github.com/ulule/limiter/v3"
	limiter_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiter_store "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	var err error

	// Set up database
	db_dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(postgres.Open(db_dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// Run migrations
	if err := DB.AutoMigrate(&Link{}); err != nil {
		panic(err.Error())
	}

	// Set up the cache
	redis_host := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr: redis_host + ":6379",
	})
	Cache = cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(2048, 32*time.Minute),
	})

	// Add middleware
	rate, _ := limiter.NewRateFromFormatted("60-S")
	store := limiter_store.NewStoreWithOptions(limiter.StoreOptions{Prefix: "limiter_gin", MaxRetry: 3})
	middleware := limiter_gin.NewMiddleware(limiter.New(store, rate))

	// Routers
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	router.Use(middleware)

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// URLs
	router.POST("/shorten", ShortenURL)
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK, "index.html", gin.H{},
		)
	})
	router.GET("/r/*susURI", ResolveURL)

	router.Run(":8080")
}
