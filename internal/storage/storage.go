package storage

import (
	"context"

	"github.com/Kapeland/task-Avito/internal/storage/db"
)

type PostgresStorage struct {
	DB *db.PgDatabase
}

func NewPostgresStorage(ctx context.Context) (PostgresStorage, error) {
	var dbStorage PostgresStorage
	database, err := db.NewPostgres(ctx)
	if err != nil {
		return PostgresStorage{}, err
	}
	dbStorage.DB = database
	return dbStorage, nil
}

func (s *PostgresStorage) Close() error {
	err := s.DB.Close()
	return err
}
