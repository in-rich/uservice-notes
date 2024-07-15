package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/samber/lo"
)

type ListNotesService interface {
	Exec(ctx context.Context, selector *models.ListNotes) ([]*models.Note, error)
}

type listNotesServiceImpl struct {
	listNotesRepository dao.ListNotesRepository
}

func (s *listNotesServiceImpl) Exec(ctx context.Context, selector *models.ListNotes) ([]*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(selector); err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	notes, err := s.listNotesRepository.ListNotes(
		ctx,
		selector.AuthorID,
		lo.Map(selector.Filters, func(item models.ListNotesFilter, index int) dao.ListNoteFilter {
			return dao.ListNoteFilter{
				PublicIdentifier: item.PublicIdentifier,
				Target:           entities.Target(item.Target),
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Note, len(notes))
	for i, note := range notes {
		result[i] = &models.Note{
			PublicIdentifier: note.PublicIdentifier,
			AuthorID:         note.AuthorID,
			Target:           string(note.Target),
			Content:          note.Content,
			UpdatedAt:        note.UpdatedAt,
		}
	}

	return result, nil
}

func NewListNotesService(listNotesRepository dao.ListNotesRepository) ListNotesService {
	return &listNotesServiceImpl{
		listNotesRepository: listNotesRepository,
	}
}
