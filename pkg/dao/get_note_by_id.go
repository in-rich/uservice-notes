package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type GetNoteByIDRepository interface {
	GetNoteByID(ctx context.Context, authorID string, noteID uuid.UUID) (*entities.Note, error)
}

type getNotesByIDRepositoryImpl struct {
	db bun.IDB
}

func (r *getNotesByIDRepositoryImpl) GetNoteByID(ctx context.Context, authorID string, noteID uuid.UUID) (*entities.Note, error) {
	note := &entities.Note{}

	err := r.db.NewSelect().
		Model(note).
		Where("author_id = ?", authorID).
		Where("id = ?", noteID).
		Order("updated_at DESC").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}

		return nil, err
	}

	return note, nil
}

func NewGetNoteByIDRepository(db bun.IDB) GetNoteByIDRepository {
	return &getNotesByIDRepositoryImpl{
		db: db,
	}
}
