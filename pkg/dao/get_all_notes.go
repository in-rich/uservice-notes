package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type GetAllNotesRepository interface {
	GetAllNotes(ctx context.Context, limit int64, offset int64) ([]*entities.Note, error)
}

type getAllNotesRepositoryImpl struct {
	db bun.IDB
}

func (r *getAllNotesRepositoryImpl) GetAllNotes(ctx context.Context, limit int64, offset int64) ([]*entities.Note, error) {
	notes := make([]*entities.Note, 0)

	err := r.db.NewSelect().
		Model(&notes).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("updated_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NewGetAllNotesRepository(db bun.IDB) GetAllNotesRepository {
	return &getAllNotesRepositoryImpl{
		db: db,
	}
}
