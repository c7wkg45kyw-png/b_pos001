package logger

import "go.uber.org/zap"

func New(env string) *zap.Logger {
	if env == "production" {
		log, _ := zap.NewProduction()
		return log
	}
	log, _ := zap.NewDevelopment()
	return log
}
