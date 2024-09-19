package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteNoteRepository interface {
	DeleteNote(ctx context.Context, author string, target entities.Target, publicIdentifier string) (*entities.Note, error)
}

type deleteNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteNoteRepositoryImpl) DeleteNote(
	ctx context.Context, author string, target entities.Target, publicIdentifier string,
) (*entities.Note, error) {
	note := &entities.Note{}

	_, err := r.db.NewDelete().
		Model(note).
		Where("author_id = ?", author).
		Where("public_identifier = ?", publicIdentifier).
		Where("target = ?", target).
		Returning("id").
		Exec(ctx)

	return note, err
}

func NewDeleteNoteRepository(db bun.IDB) DeleteNoteRepository {
	return &deleteNoteRepositoryImpl{
		db: db,
	}
}
