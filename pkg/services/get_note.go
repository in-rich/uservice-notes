package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
)

type GetNoteService interface {
	Exec(ctx context.Context, selector *models.GetNote) (*models.Note, error)
}

type getNoteServiceImpl struct {
	getNoteRepository dao.GetNoteRepository
}

func (s *getNoteServiceImpl) Exec(ctx context.Context, selector *models.GetNote) (*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(selector); err != nil {
		return nil, errors.Join(ErrInvalidNoteSelector, err)
	}

	note, err := s.getNoteRepository.GetNote(ctx, selector.AuthorID, entities.Target(selector.Target), selector.PublicIdentifier)
	if err != nil {
		return nil, err
	}

	return &models.Note{
		PublicIdentifier: note.PublicIdentifier,
		AuthorID:         note.AuthorID,
		Target:           string(note.Target),
		Content:          note.Content,
		UpdatedAt:        note.UpdatedAt,
	}, nil
}

func NewGetNoteService(getNoteRepository dao.GetNoteRepository) GetNoteService {
	return &getNoteServiceImpl{
		getNoteRepository: getNoteRepository,
	}
}
