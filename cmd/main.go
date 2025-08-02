package main

import (
	"log"

	"github.com/bryantaolong/system/internal/config"
	"github.com/bryantaolong/system/internal/router"
	"github.com/bryantaolong/system/internal/service"
	"github.com/bryantaolong/system/pkg/db"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db := db.Init(cfg)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass,
		DB:       0,
	})

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	authService := service.NewAuthService(db, redisClient)
	userService := service.NewUserService(db, authService)

	router := router.NewRouter(redisClient, authService, userService)

	log.Println("🚀 项目已启动，监听 :8080")
	log.Fatal(router.Run(":8080"))
}
