package models

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type AuthStorager interface {
	GetUserSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (structs.UserSecret, error)
	CreateUserSecret(ctx context.Context, userSecret structs.UserSecret) error
}

const validHoursNum = 24

// RegisterUser registers/auth user + always new JWT
func (m *ModelAuth) RegisterUser(ctx context.Context, info structs.RegisterUserInfo) (string, error) {
	lgr := logger.GetLogger()

	err := m.us.CreateUserST(ctx, info)
	if err != nil {
		if errors.Is(err, ErrUserConflict) {
			// Пользователь с таким логином уже есть, значит сверим пароль
			isPassCorrect, err := m.us.CheckPasswordST(ctx, structs.AuthUserInfo{Login: info.Login, Pswd: info.Pswd})
			if err != nil {
				lgr.Error(err.Error(), "ModelAuth", "LoginUser", "CheckPasswordST")

				return "", err
			}
			if !isPassCorrect { // Это значит кто-то вводит существующий логин, но другой пароль
				return "", ErrBadCredentials
			}
		} else {
			lgr.Error(err.Error(), "ModelAuth", "RegisterUser", "CreateUserDB")

			return "", err
		}
	}

	key, err := genKey(64)
	if err != nil {
		lgr.Error(err.Error(), "ModelAuth", "RegisterUser", "genKey")

		return "", err
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		lgr.Error(err.Error(), "ModelAuth", "RegisterUser", "NewV4")

		return "", err
	}

	payload := jwt.MapClaims{
		"sub": info.Login,
		"sID": sessionID.String(),
		"exp": time.Now().Add(time.Hour * validHoursNum).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokStr, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		lgr.Error(err.Error(), "ModelAuth", "RegisterUser", "SignedString")

		return "", err
	}

	userSecret := structs.UserSecret{
		Login:     info.Login,
		Secret:    key,
		SessionID: sessionID.String(),
	}
	err = m.as.CreateUserSecret(ctx, userSecret)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			// По идее такая ситуация может быть, когда секрет одинаковы, но шанс такого очень мал.
			return "", ErrConflict
		}
		lgr.Error(err.Error(), "ModelAuth", "RegisterUser", "CreateUserSecret")

		return "", err
	}

	return tokStr, nil
}

// GetUserSecretByLoginAndSession get user secret by given login and sessionID
// Returns ErrNotFound or err
func (m *ModelAuth) GetUserSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (structs.UserSecret, error) {
	userSecret, err := m.as.GetUserSecretByLoginAndSession(ctx, lgnSsn)

	return userSecret, err
}

func genKey(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		if (n >= 48 && n <= 57) || (n >= 65 && n <= 90) || (n >= 97 && n <= 122) {
			result += string(rune(n))
		}
	}
}
