package storage

import (
	"context"
	"errors"
	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/storage/repository"
)

type AuthRepo interface {
	CreateUserSecret(ctx context.Context, userSecret *structs.UserSecret) error
	GetSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (*structs.UserSecret, error)
	GetLoginBySecret(ctx context.Context, secret string) (string, error)
}

type AuthStorage struct {
	authRepo AuthRepo
}

func NewAuthStorage(authRepo AuthRepo) AuthStorage {
	return AuthStorage{authRepo: authRepo}
}

// GetUserLoginBySecret secret.
// Returns models.ErrNotFound or err
func (s *AuthStorage) GetUserLoginBySecret(ctx context.Context, secret string) (string, error) {
	login, err := s.authRepo.GetLoginBySecret(ctx, secret)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return "", models.ErrNotFound
		}
		return "", err
	}
	return login, nil
}

// GetUserSecretByLogin secret
// Returns models.ErrNotFound or err
func (s *AuthStorage) GetUserSecretByLoginAndSession(ctx context.Context, lgnSsn structs.UserSecret) (structs.UserSecret, error) {
	userSecret, err := s.authRepo.GetSecretByLoginAndSession(ctx, lgnSsn)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return structs.UserSecret{}, models.ErrNotFound
		}
		return structs.UserSecret{}, err
	}
	return *userSecret, nil
}

// CreateUserSecret secret
// Returns models.ErrConflict or err
func (s *AuthStorage) CreateUserSecret(ctx context.Context, userSecret structs.UserSecret) error {
	err := s.authRepo.CreateUserSecret(ctx, &userSecret)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return models.ErrConflict
		}
		return err
	}
	return nil
}
