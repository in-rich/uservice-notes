package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type ListNotesByAuthorRepository interface {
	ListNotesByAuthor(ctx context.Context, authorID string, limit int, offset int) ([]*entities.Note, error)
}

type listNotesByAuthorRepositoryImpl struct {
	db bun.IDB
}

func (r *listNotesByAuthorRepositoryImpl) ListNotesByAuthor(ctx context.Context, authorID string, limit int, offset int) ([]*entities.Note, error) {
	notes := make([]*entities.Note, 0)

	err := r.db.NewSelect().
		Model(&notes).
		Where("author_id = ?", authorID).
		Limit(limit).
		Offset(offset).
		Order("updated_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NewListNotesByAuthorRepository(db bun.IDB) ListNotesByAuthorRepository {
	return &listNotesByAuthorRepositoryImpl{
		db: db,
	}
}
