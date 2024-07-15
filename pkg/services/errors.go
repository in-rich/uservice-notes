package services

import "errors"

var (
	ErrInvalidNoteSelector     = errors.New("invalid note selector")
	ErrInvalidNoteUpdate       = errors.New("invalid note update")
	ErrNotesUpdateLimitReached = errors.New("notes update limit reached")
)
