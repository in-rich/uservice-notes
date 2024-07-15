package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteNoteRepository interface {
	DeleteNote(ctx context.Context, author string, target entities.Target, publicIdentifier string) error
}

type deleteNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteNoteRepositoryImpl) DeleteNote(
	ctx context.Context, author string, target entities.Target, publicIdentifier string,
) error {
	_, err := r.db.NewDelete().
		Model(&entities.Note{}).
		Where("author_id = ?", author).
		Where("public_identifier = ?", publicIdentifier).
		Where("target = ?", target).
		Exec(ctx)

	return err
}

func NewDeleteNoteRepository(db bun.IDB) DeleteNoteRepository {
	return &deleteNoteRepositoryImpl{
		db: db,
	}
}
