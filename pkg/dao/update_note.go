package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"time"
)

type UpdateNoteData struct {
	Content string
}

type UpdateNoteRepository interface {
	UpdateNote(ctx context.Context, author string, target entities.Target, publicIdentifier string, data *UpdateNoteData) (*entities.Note, error)
}

type updateNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *updateNoteRepositoryImpl) UpdateNote(
	ctx context.Context, author string, target entities.Target, publicIdentifier string, data *UpdateNoteData,
) (*entities.Note, error) {
	note := &entities.Note{
		Content:   data.Content,
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	res, err := r.db.NewUpdate().
		Model(note).
		Column("content", "updated_at").
		Where("author_id = ?", author).
		Where("public_identifier = ?", publicIdentifier).
		Where("target = ?", target).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrNoteNotFound
	}

	return note, nil
}

func NewUpdateNoteRepository(db bun.IDB) UpdateNoteRepository {
	return &updateNoteRepositoryImpl{
		db: db,
	}
}
