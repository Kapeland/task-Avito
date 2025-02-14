package models

import (
	"context"

	"github.com/Kapeland/task-Avito/internal/models/structs"
)

type ModelUsers struct {
	us UsersStorager
}

type ModelAuth struct {
	as AuthStorager
	us UsersStorager
}

func NewModelUsers(us UsersStorager) ModelUsers {
	return ModelUsers{us}
}
func NewModelAuth(as AuthStorager, us UsersStorager) ModelAuth {
	return ModelAuth{as, us}
}

type AuthModelManager interface {
	RegisterUser(ctx context.Context, info structs.RegisterUserInfo) (string, error)
	GetUserSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (structs.UserSecret, error)
}

type UsersModelManager interface {
	SendCoin(ctx context.Context, operation structs.SendCoinInfo) error
	BuyItem(ctx context.Context, item string, login string) error
	Info(ctx context.Context, login string) (structs.AccInfo, error)
}
