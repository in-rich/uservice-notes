package dao

import "errors"

var (
	ErrNoteAlreadyExists = errors.New("note already exists")
	ErrNoteNotFound      = errors.New("note not found")
)
