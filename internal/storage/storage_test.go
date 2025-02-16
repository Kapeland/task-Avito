package storage

import (
	"context"
	"testing"

	"github.com/Kapeland/task-Avito/internal/storage/db"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
)

func TestPostgresStorage_Close(t *testing.T) {
	type fields struct {
		DB *db.PgDatabase
	}

	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		t.Error(err)
	}
	cfg := config.GetConfig()
	logger.CreateLogger(&cfg)
	tmpDB, err := db.NewPostgres(ctx)

	if err != nil {
		t.Error("NewPostgres: " + err.Error())
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Success close",
			fields:  fields{DB: tmpDB},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostgresStorage{
				DB: tt.fields.DB,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
