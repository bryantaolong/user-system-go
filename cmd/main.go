package main

import (
	"log"

	"github.com/bryantaolong/system/internal/config"
	"github.com/bryantaolong/system/internal/router"
	"github.com/bryantaolong/system/internal/service"
	"github.com/bryantaolong/system/internal/service/redis"
	"github.com/bryantaolong/system/pkg/db"
	goredis "github.com/go-redis/redis/v8" // ✅ 官方库
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db := db.Init(cfg)

	// ✅ 使用官方库
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass,
		DB:       0,
	})

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	rdb := redis.NewRedisStringService(redisClient, logger)
	authService := service.NewAuthService(db, rdb, cfg.JWTSecret)
	userService := service.NewUserService(db, authService)

	router := router.NewRouter(db, authService, userService)

	log.Println("🚀 项目已启动，监听 :8080")
	log.Fatal(router.Run(":8080"))
}
