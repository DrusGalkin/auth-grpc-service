package storage

import "errors"

var (
	ErrUserExist    = errors.New("Пользователь уже существует")
	ErrUserNotFound = errors.New("Пользователь не найден")
)
