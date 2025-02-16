package users

import (
	"context"
	"database/sql"
	"errors"

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

// CreateUserDB create user
func (r *Repo) CreateUserDB(ctx context.Context, info structs.RegisterUserInfo) error {
	lgr := logger.GetLogger()

	id := 0

	tx, err := r.db.(*db.PgDatabase).BeginX(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.users(login, password_hash)
				VALUES($1, crypt($2, gen_salt('bf'))) returning id;`, info.Login, info.Pswd).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return repository.ErrDuplicateKey
		}
		lgr.Error(err.Error(), "Repo", "CreateUserDB", "INSERT1")
		return err
	}
	tmp := ""
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.account(login)
				VALUES($1) returning login;`, info.Login).Scan(&tmp)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return repository.ErrDuplicateKey
		}
		lgr.Error(err.Error(), "Repo", "CreateUserDB", "INSERT2")

		return err
	}

	if err := tx.Commit(); err != nil {
		lgr.Error(err.Error(), "Repo", "CreateUserDB", "Commit")
		return err
	}

	return nil
}

// VerifyPasswordDB checks whether the password is correct or no.
func (r *Repo) VerifyPasswordDB(ctx context.Context, info structs.AuthUserInfo) (bool, error) {
	lgr := logger.GetLogger()

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
			lgr.Error(err.Error(), "Repo", "VerifyPasswordDB", "Commit")

			return false, err
		}
		return false, nil
	case err != nil:
		lgr.Error(err.Error(), "Repo", "VerifyPasswordDB", "SELECT")
		return false, err
	default:
		if err := tx.Commit(); err != nil {
			lgr.Error(err.Error(), "Repo", "VerifyPasswordDB", "def-Commit")
			return false, err
		}
		return isValid, nil
	}
}

// SendCoinDB send coin to user
func (r *Repo) SendCoinDB(ctx context.Context, operation structs.SendCoinInfo) error {
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
			// По идее такой ошибки быть не может, так как мы ранее JWT проверили и его владелец точно есть
			return repository.ErrObjectNotFound
		}

		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23514" {
			// check constraint from PostgreSQL
			// Проверяет balance
			return repository.ErrCheckConstraint
		}

		lgr.Error(err.Error(), "Repo", "SendCoinDB", "UPDATE1")

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
		lgr.Error(err.Error(), "Repo", "SendCoinDB", "UPDATE2")

		return err
	}
	// Сохранили операцию
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.user_operations(sender, recipient, amount)
				VALUES($1, $2, $3) returning sender;`, operation.From, operation.To, operation.Amount).Scan(&tmp)

	if err != nil {
		lgr.Error(err.Error(), "Repo", "SendCoinDB", "INSERT")

		return err
	}

	if err := tx.Commit(); err != nil {
		lgr.Error(err.Error(), "Repo", "SendCoinDB", "Commit")

		return err
	}

	return nil
}

// BuyItemDB buy item
func (r *Repo) BuyItemDB(ctx context.Context, item string, login string) error {
	lgr := logger.GetLogger()

	_, ok := items[item]

	if !ok {
		return repository.ErrNoSuchItem
	}

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
		lgr.Error(err.Error(), "Repo", "BuyItemDB", "UPDATE")

		return err
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO users_schema.user_items(login, item)
				VALUES($1, $2) returning login;`, login, item).Scan(&tmp)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code != "23505" { // Проверка, что это не ошибка дублирования, так как это не проблема
			lgr.Error(err.Error(), "Repo", "BuyItemDB", "INSERT")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		lgr.Error(err.Error(), "Repo", "BuyItemDB", "Commit")
		return err
	}

	return nil
}

// GetInfoDB get all info about user
func (r *Repo) GetInfoDB(ctx context.Context, login string) (*structs.AccInfo, error) {
	lgr := logger.GetLogger()

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
			return &structs.AccInfo{}, repository.ErrObjectNotFound
		}

		lgr.Error(err.Error(), "Repo", "GetInfoDB", "SELECT1")

		return &structs.AccInfo{}, err
	}

	err = tx.SelectContext(ctx, &accInfo.Inventory,
		`SELECT item, COUNT(*) AS cnt FROM users_schema.user_items WHERE login=$1 GROUP BY item;`, login)
	if err != nil {
		lgr.Error(err.Error(), "Repo", "GetInfoDB", "SELECT2")

		return &structs.AccInfo{}, err
	}

	err = tx.SelectContext(ctx, &accInfo.CoinHistory.Received,
		`SELECT sender, SUM(amount) AS amount FROM users_schema.user_operations WHERE recipient=$1 GROUP BY sender;`, login)
	if err != nil {
		lgr.Error(err.Error(), "Repo", "GetInfoDB", "SELECT3")

		return &structs.AccInfo{}, err
	}

	err = tx.SelectContext(ctx, &accInfo.CoinHistory.Sent,
		`SELECT recipient, SUM(amount) AS amount FROM users_schema.user_operations WHERE sender=$1 GROUP BY recipient;`, login)
	if err != nil {
		lgr.Error(err.Error(), "Repo", "GetInfoDB", "SELECT4")

		return &structs.AccInfo{}, err
	}

	if err := tx.Commit(); err != nil {
		lgr.Error(err.Error(), "Repo", "GetInfoDB", "Commit")

		return &structs.AccInfo{}, err
	}

	return &accInfo, nil
}
