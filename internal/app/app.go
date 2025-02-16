package app

import (
	"context"
	"flag"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/services"
	"github.com/Kapeland/task-Avito/internal/storage"
	"github.com/Kapeland/task-Avito/internal/storage/repository/postgresql/auth"
	"github.com/Kapeland/task-Avito/internal/storage/repository/postgresql/users"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/pressly/goose/v3"
)

func Start(cfg *config.Config, lgr *logger.Logger) error {
	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	ctx := context.Background()
	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		lgr.Error(err.Error(), "App", "Start", "NewPostgresStorage")
		return err
	}
	defer dbStor.Close()

	if *migration {
		if err := goose.Up(dbStor.DB.GetDB().DB, cfg.Database.Migrations); err != nil {
			lgr.Error("Migration failed: "+err.Error(), "App", "Start", " goose.Up")

			return err
		}
	}

	usersRepo := users.New(dbStor.DB)
	authRepo := auth.New(dbStor.DB)

	authStorage := storage.NewAuthStorage(authRepo)
	usersStorage := storage.NewUsersStorage(usersRepo)

	umdl := models.NewModelUsers(&usersStorage)
	amdl := models.NewModelAuth(&authStorage, &usersStorage)

	serv := services.NewService(&umdl, &amdl)

	return serv.Launch(cfg, lgr)
}
