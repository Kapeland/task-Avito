package users

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/storage/db"
	"github.com/Kapeland/task-Avito/internal/storage/repository"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	db db.DBops
}

func New(db db.DBops) *Repo {
	return &Repo{db: db}
}

// CreateUser create user
func (r *Repo) CreateUser(ctx context.Context, info structs.RegisterUserInfo) (int, error) {
	id := 0

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.users(login, password_hash)
				VALUES($1, crypt($2, gen_salt('bf'))) returning id;`, info.Login, info.Pswd).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return id, repository.ErrDuplicateKey
		}
		return id, err
	}
	tmp := ""
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.account(login)
				VALUES($1) returning login;`, info.Login).Scan(&tmp)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return id, repository.ErrDuplicateKey
		}
		return id, err
	}

	if err := tx.Commit(); err != nil {
		slog.Info(repository.ErrContextClosed.Error())
		slog.Error(err.Error())
		return id, err
	}

	return id, nil
}

// VerifyPassword checks whether the password is correct or no.
func (r *Repo) VerifyPassword(ctx context.Context, info structs.AuthUserInfo) (bool, error) {
	isValid := false

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`SELECT (password_hash = crypt($1, password_hash)) 
    			AS password_match 
				FROM users_schema.users
				WHERE login = $2 ;`, info.Pswd, info.Login).Scan(&isValid)

	switch {
	case err != nil && (errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows)):
		if err := tx.Commit(); err != nil {
			slog.Info(repository.ErrContextClosed.Error())
			slog.Error(err.Error())
			return false, err
		}
		return false, nil
	case err != nil:
		return false, err
	default:
		if err := tx.Commit(); err != nil {
			slog.Info(repository.ErrContextClosed.Error())
			slog.Error(err.Error())
			return false, err
		}
		return isValid, nil
	}
}

// GetUserByLogin get user
func (r *Repo) GetUserByLogin(ctx context.Context, login string) (*structs.User, error) {
	var info structs.User

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.GetContext(ctx, &info,
		`SELECT id, login, password_hash FROM users_schema.users WHERE login=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		slog.Info(repository.ErrContextClosed.Error())
		slog.Error(err.Error())
		return nil, err
	}

	return &info, nil
}

// SendCoinTo send coin to user
func (r *Repo) SendCoinTo(ctx context.Context, operation structs.SendCoinInfo) error {
	lgr := logger.GetLogger()

	tmp := ""

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вычли у отправителя
	err = tx.QueryRowContext(ctx,
		`UPDATE users_schema.account SET balance = balance- $1
				WHERE login = $2 returning login;`, operation.Amount, operation.From).Scan(&tmp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		lgr.Error(err.Error(), "Repo", "SendCoinTo", "UPDATE")

		return err
	}

	// Добавили получателю
	err = tx.QueryRowContext(ctx,
		`UPDATE users_schema.account SET balance = balance + $1
				WHERE login = $2 returning login;`, operation.Amount, operation.To).Scan(&tmp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		lgr.Error(err.Error(), "Repo", "SendCoinTo", "UPDATE")

		return err
	}
	// Сохранили операцию
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.user_operations(sender, recipient, amount)
				VALUES($1, $2, $3) returning sender;`, operation.From, operation.To, operation.Amount).Scan(&tmp)

	if err != nil {
		lgr.Error(err.Error(), "Repo", "SendCoinTo", "INSERT")

		return err
	}

	if err := tx.Commit(); err != nil {
		slog.Info(repository.ErrContextClosed.Error())
		slog.Error(err.Error())
		return err
	}

	return nil
}

// BuyItem buy item
func (r *Repo) BuyItem(ctx context.Context, item string, login string) error {
	lgr := logger.GetLogger()

	tmp := ""

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`UPDATE users_schema.account SET balance = balance- $1
				WHERE login = $2 returning login;`, items[item], login).Scan(&tmp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		lgr.Error(err.Error(), "Repo", "BuyItem", "UPDATE")

		return err
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.user_items(login, item)
				VALUES($1, $2) returning login;`, login, item).Scan(&tmp)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code != "23505" { // Проверка, что это не ошибка дублирования, так как это не проблема
			lgr.Error(err.Error(), "Repo", "BuyItem", "INSERT")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Info(repository.ErrContextClosed.Error())
		slog.Error(err.Error())
		return err
	}

	return nil
}

// GetInfo get all info about user
func (r *Repo) GetInfo(ctx context.Context, login string) (*structs.AccInfo, error) {
	var accInfo structs.AccInfo

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return &structs.AccInfo{}, err
	}
	defer tx.Rollback()

	err = tx.GetContext(ctx, &accInfo.Coins,
		`SELECT balance FROM users_schema.account WHERE login=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return &structs.AccInfo{}, err
	}

	err = tx.SelectContext(ctx, &accInfo.Inventory,
		`SELECT item, COUNT(*) AS cnt FROM users_schema.user_items WHERE login=$1 GROUP BY item;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return &structs.AccInfo{}, err
	}

	err = tx.SelectContext(ctx, &accInfo.CoinHistory.Received,
		`SELECT sender, amount FROM users_schema.user_operations WHERE recipient=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return &structs.AccInfo{}, err
	}
	err = tx.SelectContext(ctx, &accInfo.CoinHistory.Sent,
		`SELECT recipient, amount FROM users_schema.user_operations WHERE sender=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return &structs.AccInfo{}, err
	}

	if err := tx.Commit(); err != nil {
		slog.Info(repository.ErrContextClosed.Error())
		slog.Error(err.Error())
		return &structs.AccInfo{}, err
	}

	return &accInfo, nil
}

//TODO: проверить, что все ошибки правильно возвращаются
//TODO: добавить логи или мб удалить ваще
