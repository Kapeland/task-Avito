package repository

import "errors"

var ErrObjectNotFound = errors.New("not found")

var ErrDuplicateKey = errors.New("duplicate key")

var ErrNoSuchItem = errors.New("unknown item")

var ErrCheckConstraint = errors.New("violating check constraint")
