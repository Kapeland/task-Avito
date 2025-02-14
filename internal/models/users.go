package models

import (
	"context"

	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
)

type UsersStorager interface {
	CreateUser(ctx context.Context, info structs.RegisterUserInfo) (int, error)
	GetUserByLogin(ctx context.Context, login string) (structs.User, error)
	CheckPassword(ctx context.Context, info structs.AuthUserInfo) (bool, error)
	SendCoinTo(ctx context.Context, operation structs.SendCoinInfo) error
	BuyItem(ctx context.Context, item string, login string) error
	GetInfo(ctx context.Context, login string) (structs.AccInfo, error)
}

func (m *ModelUsers) SendCoin(ctx context.Context, operation structs.SendCoinInfo) error {
	lgr := logger.GetLogger()

	err := m.us.SendCoinTo(ctx, operation)

	if err != nil {
		lgr.Error(err.Error(), "ModelUsers", "SendCoin", "SendCoinTo")

		return err
	}

	return nil
}
func (m *ModelUsers) BuyItem(ctx context.Context, item string, login string) error {
	lgr := logger.GetLogger()

	err := m.us.BuyItem(ctx, item, login)

	if err != nil {
		lgr.Error(err.Error(), "ModelUsers", "BuyItem", "BuyItem")

		return err
	}

	return nil
}

func (m *ModelUsers) Info(ctx context.Context, login string) (structs.AccInfo, error) {
	lgr := logger.GetLogger()

	accInfo, err := m.us.GetInfo(ctx, login)

	if err != nil {
		lgr.Error(err.Error(), "ModelUsers", "Info", "Info")

		return structs.AccInfo{}, err
	}

	return accInfo, nil
}
