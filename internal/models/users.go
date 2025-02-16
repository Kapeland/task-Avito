package models

import (
	"context"
	"errors"

	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
)

type UsersStorager interface {
	CreateUserST(ctx context.Context, info structs.RegisterUserInfo) error
	CheckPasswordST(ctx context.Context, info structs.AuthUserInfo) (bool, error)
	SendCoinST(ctx context.Context, operation structs.SendCoinInfo) error
	BuyItemST(ctx context.Context, item string, login string) error
	GetInfoST(ctx context.Context, login string) (structs.AccInfo, error)
}

func (m *ModelUsers) SendCoin(ctx context.Context, operation structs.SendCoinInfo) error {
	lgr := logger.GetLogger()

	err := m.us.SendCoinST(ctx, operation)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrUserNotFound
		}
		if errors.Is(err, ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}

		lgr.Error(err.Error(), "ModelUsers", "SendCoin", "SendCoinDB")

		return err
	}

	return nil
}

func (m *ModelUsers) BuyItem(ctx context.Context, item string, login string) error {
	lgr := logger.GetLogger()

	err := m.us.BuyItemST(ctx, item, login)

	if err != nil {
		if errors.Is(err, ErrNoSuchItem) {
			return ErrNoSuchItem
		}

		lgr.Error(err.Error(), "ModelUsers", "BuyItemDB", "BuyItemDB")

		return err
	}

	return nil
}

func (m *ModelUsers) Info(ctx context.Context, login string) (structs.AccInfo, error) {
	lgr := logger.GetLogger()

	accInfo, err := m.us.GetInfoST(ctx, login)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return structs.AccInfo{}, ErrUserNotFound
		}

		lgr.Error(err.Error(), "ModelUsers", "Info", "Info")

		return structs.AccInfo{}, err
	}

	return accInfo, nil
}
