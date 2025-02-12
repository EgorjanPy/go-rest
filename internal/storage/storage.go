package storage

import "errors"

var (
	ErrURLExists   = errors.New("error url ulready exists")
	ErrURLNotFound = errors.New("error url not found")
)
