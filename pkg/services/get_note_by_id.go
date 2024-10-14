package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/models"
)

type GetNoteByIDService interface {
	Exec(ctx context.Context, selector *models.GetNoteByID) (*models.Note, error)
}

type getNoteByIDServiceImpl struct {
	getNoteByIDRepository dao.GetNoteByIDRepository
}

func (s *getNoteByIDServiceImpl) Exec(ctx context.Context, selector *models.GetNoteByID) (*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(selector); err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	noteID, err := uuid.Parse(selector.NoteID)
	if err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	note, err := s.getNoteByIDRepository.GetNoteByID(ctx, selector.AuthorID, noteID)
	if err != nil {
		return nil, err
	}

	return &models.Note{
		ID:               note.ID.String(),
		PublicIdentifier: note.PublicIdentifier,
		AuthorID:         note.AuthorID,
		Target:           string(note.Target),
		Content:          note.Content,
		UpdatedAt:        note.UpdatedAt,
	}, nil
}

func NewGetNoteByIDService(getNoteByIDRepository dao.GetNoteByIDRepository) GetNoteByIDService {
	return &getNoteByIDServiceImpl{
		getNoteByIDRepository: getNoteByIDRepository,
	}
}
