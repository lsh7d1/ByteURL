package dao

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrExistData = errors.New("existing data")
)
