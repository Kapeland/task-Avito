package models

import "errors"

var ErrNotFound = errors.New("item not found")

var ErrConflict = errors.New("item already exists")

var ErrUserConflict = errors.New("user already exists")

var ErrBadCredentials = errors.New("bad password or login")

var ErrInvalidToken = errors.New("invalid token")

var ErrNoSuchItem = errors.New("unknown item")

var ErrUserNotFound = errors.New("user not found")

var ErrInsufficientBalance = errors.New("insufficient balance")
