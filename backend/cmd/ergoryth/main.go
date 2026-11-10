package main

import "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/logger"

func main() {
	log := logger.NewLogger(logger.NewConfigMust())

	log.Info("application started")
}
