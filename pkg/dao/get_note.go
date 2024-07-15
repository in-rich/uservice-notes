package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/uptrace/bun"
)

type GetNoteRepository interface {
	GetNote(ctx context.Context, author string, target entities.Target, publicIdentifier string) (*entities.Note, error)
}

type getNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *getNoteRepositoryImpl) GetNote(
	ctx context.Context, author string, target entities.Target, publicIdentifier string,
) (*entities.Note, error) {
	note := new(entities.Note)

	err := r.db.NewSelect().
		Model(note).
		Where("author_id = ?", author).
		Where("public_identifier = ?", publicIdentifier).
		Where("target = ?", target).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}

		return nil, err
	}

	return note, nil
}

func NewGetNoteRepository(db bun.IDB) GetNoteRepository {
	return &getNoteRepositoryImpl{
		db: db,
	}
}
