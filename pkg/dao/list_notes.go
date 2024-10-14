package dao

import (
	"context"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type ListNoteFilter struct {
	PublicIdentifier string
	Target           entities.Target
}

type ListNotesRepository interface {
	ListNotes(ctx context.Context, authorID string, filters []ListNoteFilter) ([]*entities.Note, error)
}

type listNotesRepositoryImpl struct {
	db bun.IDB
}

func (r *listNotesRepositoryImpl) ListNotes(ctx context.Context, authorID string, filters []ListNoteFilter) ([]*entities.Note, error) {
	notes := make([]*entities.Note, 0)

	userNotesFilter := lo.Reduce(filters, func(agg []string, item ListNoteFilter, index int) []string {
		if item.Target == entities.TargetCompany {
			return agg
		}
		return append(agg, item.PublicIdentifier)
	}, []string{})

	companyNotesFilter := lo.Reduce(filters, func(agg []string, item ListNoteFilter, index int) []string {
		if item.Target == entities.TargetUser {
			return agg
		}
		return append(agg, item.PublicIdentifier)
	}, []string{})

	usersQuery := r.db.NewSelect().
		Model(&notes).
		Where("author_id = ?", authorID).
		Where("target = ?", entities.TargetUser).
		Where("public_identifier IN (?)", bun.In(userNotesFilter))

	companiesQuery := r.db.NewSelect().
		Model(&notes).
		Where("author_id = ?", authorID).
		Where("target = ?", entities.TargetCompany).
		Where("public_identifier IN (?)", bun.In(companyNotesFilter))

	err := usersQuery.Union(companiesQuery).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NewListNotesRepository(db bun.IDB) ListNotesRepository {
	return &listNotesRepositoryImpl{
		db: db,
	}
}
