package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"user-restriction-manager/internal/api"
	"user-restriction-manager/internal/api/handler"
	"user-restriction-manager/internal/config"
	"user-restriction-manager/internal/core/service"
	"user-restriction-manager/internal/repository/postgres"
	redisRepo "user-restriction-manager/internal/repository/redis"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Database connection
	dbConnStr := cfg.DBConnectionString()
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer rdb.Close()
	
	// Repositories
	restrictionRepo := postgres.NewRestrictionRepository(db)
	appealRepo := postgres.NewAppealRepository(db)
	cacheRepo := redisRepo.NewCacheRepository(rdb)

	// Services
	restrictionService := service.NewRestrictionService(restrictionRepo, cacheRepo)
	appealService := service.NewAppealService(appealRepo, restrictionRepo)

	// Handlers
	restrictionHandler := handler.NewRestrictionHandler(restrictionService)
	appealHandler := handler.NewAppealHandler(appealService)

	// Router
	r := api.NewRouter(restrictionHandler, appealHandler)

	// Server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s (%s)", addr, cfg.AppEnv)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
