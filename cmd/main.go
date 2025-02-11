package main

import (
	"github.com/Kapeland/task-Avito/internal/app"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"log/slog"
	"os"
)

func main() {
	if err := config.ReadConfigYAML(); err != nil {
		slog.Error("Failed init configuration")
		os.Exit(1)
	}
	cfg := config.GetConfig()
	lgr := logger.CreateLogger(&cfg)

	lgr.Info("app started", "main", "", "")
	err := app.Start(&cfg, &lgr)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	lgr.Info("app finished", "main", "", "")
}
