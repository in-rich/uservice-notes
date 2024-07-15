package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/models"
)

type ListNotesByAuthorService interface {
	Exec(ctx context.Context, selector *models.ListNotesByAuthor) ([]*models.Note, error)
}

type listNotesByAuthorServiceImpl struct {
	listNotesByAuthorRepository dao.ListNotesByAuthorRepository
}

func (s *listNotesByAuthorServiceImpl) Exec(ctx context.Context, selector *models.ListNotesByAuthor) ([]*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(selector); err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	notes, err := s.listNotesByAuthorRepository.ListNotesByAuthor(ctx, selector.AuthorID, selector.Limit, selector.Offset)
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

func NewListNotesByAuthorService(listNotesByAuthorRepository dao.ListNotesByAuthorRepository) ListNotesByAuthorService {
	return &listNotesByAuthorServiceImpl{
		listNotesByAuthorRepository: listNotesByAuthorRepository,
	}
}
