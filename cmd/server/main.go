package main

import (
	"bpos001/backend/internal/config"
	"bpos001/backend/internal/database"
	"bpos001/backend/internal/handler"
	"bpos001/backend/internal/logger"
	"bpos001/backend/internal/repository"
	"bpos001/backend/internal/route"
	"bpos001/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)
	defer log.Sync()
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("migrate database", zap.Error(err))
	}
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	repo := repository.New(db)
	uc := usecase.New(repo)
	route.Register(r, cfg, route.Handlers{POS: handler.NewPOSHandler(uc)})
	log.Info("BPOS001 backend started", zap.String("port", cfg.AppPort))
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("server stopped", zap.Error(err))
	}
}
