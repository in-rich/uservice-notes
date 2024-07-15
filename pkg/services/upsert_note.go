package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
)

type UpsertNoteService interface {
	Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, error)
}

type upsertNoteServiceImpl struct {
	updateNoteRepository dao.UpdateNoteRepository
	createNoteRepository dao.CreateNoteRepository
	deleteNoteRepository dao.DeleteNoteRepository
}

func (s *upsertNoteServiceImpl) Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(note); err != nil {
		return nil, errors.Join(ErrInvalidNoteUpdate, err)
	}

	// Delete note if content is empty.
	if note.Content == "" {
		return nil, s.deleteNoteRepository.DeleteNote(ctx, note.AuthorID, entities.Target(note.Target), note.PublicIdentifier)
	}

	// Attempt to create a note.
	createdNote, err := s.createNoteRepository.CreateNote(
		ctx,
		note.AuthorID,
		entities.Target(note.Target),
		note.PublicIdentifier,
		&dao.CreateNoteData{Content: note.Content},
	)

	// Note was successfully created.
	if err == nil {
		return &models.Note{
			PublicIdentifier: createdNote.PublicIdentifier,
			AuthorID:         createdNote.AuthorID,
			Target:           string(createdNote.Target),
			Content:          createdNote.Content,
			UpdatedAt:        createdNote.UpdatedAt,
		}, nil
	}

	if !errors.Is(err, dao.ErrNoteAlreadyExists) {
		return nil, err
	}

	// Note already exists. Update it.
	updatedNote, err := s.updateNoteRepository.UpdateNote(
		ctx,
		note.AuthorID,
		entities.Target(note.Target),
		note.PublicIdentifier,
		&dao.UpdateNoteData{Content: note.Content},
	)

	if err != nil {
		return nil, err
	}

	return &models.Note{
		PublicIdentifier: updatedNote.PublicIdentifier,
		AuthorID:         updatedNote.AuthorID,
		Target:           string(updatedNote.Target),
		Content:          updatedNote.Content,
		UpdatedAt:        updatedNote.UpdatedAt,
	}, nil
}

func NewUpsertNoteService(
	updateNoteRepository dao.UpdateNoteRepository,
	createNoteRepository dao.CreateNoteRepository,
	deleteNoteRepository dao.DeleteNoteRepository,
) UpsertNoteService {
	return &upsertNoteServiceImpl{
		updateNoteRepository: updateNoteRepository,
		createNoteRepository: createNoteRepository,
		deleteNoteRepository: deleteNoteRepository,
	}
}
