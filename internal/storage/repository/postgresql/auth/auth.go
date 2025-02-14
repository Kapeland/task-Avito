package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/storage/db"
	"github.com/Kapeland/task-Avito/internal/storage/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

type Repo struct {
	db db.DBops
}

func New(db db.DBops) *Repo {
	return &Repo{db: db}
}

// CreateUserSecret create user secret
func (r *Repo) CreateUserSecret(ctx context.Context, userSecret *structs.UserSecret) error {
	tmpLgn := ""

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO auth_schema.users_secrets(login, secret, session_id)
				VALUES($1, $2, $3) returning login;`, userSecret.Login, userSecret.Secret, userSecret.SessionID).Scan(&tmpLgn)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return repository.ErrDuplicateKey
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Info("Looks like the context has been closed")
		slog.Error(err.Error())
		return err
	}

	return nil
}

// GetSecretByLoginAndSession get secret
// Returns repository.ErrObjectNotFound or err
func (r *Repo) GetSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (*structs.UserSecret, error) {
	userSecret := structs.UserSecret{}

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.GetContext(ctx, &userSecret,
		`SELECT login, secret, session_id FROM auth_schema.users_secrets WHERE login=$1 and session_id=$2;`, lgnSsn.Login, lgnSsn.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		slog.Info("Looks like the context has been closed")
		slog.Error(err.Error())
		return nil, err
	}

	return &userSecret, nil
}
