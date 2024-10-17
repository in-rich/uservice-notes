package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/models"
)

type GetAllNotesService interface {
	Exec(ctx context.Context, selector *models.GetAllNotes) ([]*models.Note, error)
}

type getAllNotesServiceImpl struct {
	getAllNotesRepository dao.GetAllNotesRepository
}

func (s *getAllNotesServiceImpl) Exec(ctx context.Context, selector *models.GetAllNotes) ([]*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(selector); err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	notes, err := s.getAllNotesRepository.GetAllNotes(ctx, selector.Limit, selector.Offset)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Note, len(notes))
	for i, note := range notes {
		result[i] = &models.Note{
			ID:               note.ID.String(),
			PublicIdentifier: note.PublicIdentifier,
			AuthorID:         note.AuthorID,
			Target:           string(note.Target),
			Content:          note.Content,
			UpdatedAt:        note.UpdatedAt,
		}
	}

	return result, nil
}

func NewGetAllNotesService(getAllNotesRepository dao.GetAllNotesRepository) GetAllNotesService {
	return &getAllNotesServiceImpl{
		getAllNotesRepository: getAllNotesRepository,
	}
}
