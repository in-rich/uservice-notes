package dao

import (
	"context"
	"errors"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type CreateNoteData struct {
	Content   string
	UpdatedAt *time.Time
}

type CreateNoteRepository interface {
	CreateNote(ctx context.Context, author string, target entities.Target, publicIdentifier string, data *CreateNoteData) (*entities.Note, error)
}

type createNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *createNoteRepositoryImpl) CreateNote(
	ctx context.Context, author string, target entities.Target, publicIdentifier string, data *CreateNoteData,
) (*entities.Note, error) {
	note := &entities.Note{
		PublicIdentifier: publicIdentifier,
		AuthorID:         author,
		Target:           target,
		Content:          data.Content,
		// This value is optional, and will not be set if it is nil.
		UpdatedAt: lo.CoalesceOrEmpty(data.UpdatedAt, lo.ToPtr(time.Now())),
	}

	if _, err := r.db.NewInsert().Model(note).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrNoteAlreadyExists
		}

		return nil, err
	}

	return note, nil
}

func NewCreateNoteRepository(db bun.IDB) CreateNoteRepository {
	return &createNoteRepositoryImpl{
		db: db,
	}
}
