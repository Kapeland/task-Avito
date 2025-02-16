package storage

import (
	"context"
	"errors"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/storage/repository"
)

type UsersRepo interface {
	CreateUserDB(ctx context.Context, info structs.RegisterUserInfo) error
	VerifyPasswordDB(ctx context.Context, info structs.AuthUserInfo) (bool, error)
	SendCoinDB(ctx context.Context, operation structs.SendCoinInfo) error
	BuyItemDB(ctx context.Context, item string, login string) error
	GetInfoDB(ctx context.Context, login string) (*structs.AccInfo, error)
}

type UsersStorage struct {
	usersRepo UsersRepo
}

func NewUsersStorage(usersRepo UsersRepo) UsersStorage {
	return UsersStorage{usersRepo: usersRepo}
}

// CreateUserST user
func (s *UsersStorage) CreateUserST(ctx context.Context, info structs.RegisterUserInfo) error {
	err := s.usersRepo.CreateUserDB(ctx, info)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return models.ErrUserConflict
		}
		return err
	}

	return nil
}

// CheckPasswordST user
func (s *UsersStorage) CheckPasswordST(ctx context.Context, info structs.AuthUserInfo) (bool, error) {
	ok, err := s.usersRepo.VerifyPasswordDB(ctx, info)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// SendCoinST user
func (s *UsersStorage) SendCoinST(ctx context.Context, operation structs.SendCoinInfo) error {
	err := s.usersRepo.SendCoinDB(ctx, operation)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return models.ErrUserNotFound
		}
		if errors.Is(err, repository.ErrCheckConstraint) {
			return models.ErrInsufficientBalance
		}
		return err
	}
	return nil
}

// BuyItemST user
func (s *UsersStorage) BuyItemST(ctx context.Context, item string, login string) error {
	err := s.usersRepo.BuyItemDB(ctx, item, login)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchItem) {
			return models.ErrNoSuchItem
		}
		return err
	}
	return nil
}

// GetInfoST user
func (s *UsersStorage) GetInfoST(ctx context.Context, login string) (structs.AccInfo, error) {
	info, err := s.usersRepo.GetInfoDB(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return structs.AccInfo{}, models.ErrUserNotFound
		}
		return structs.AccInfo{}, err
	}
	return *info, err
}
