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
	Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, string, error)
}

type upsertNoteServiceImpl struct {
	updateNoteRepository dao.UpdateNoteRepository
	createNoteRepository dao.CreateNoteRepository
	deleteNoteRepository dao.DeleteNoteRepository
}

func (s *upsertNoteServiceImpl) Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, string, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(note); err != nil {
		return nil, "", errors.Join(ErrInvalidNoteUpdate, err)
	}

	// Delete note if content is empty.
	if note.Content == "" {
		deletedNote, err := s.deleteNoteRepository.DeleteNote(ctx, note.AuthorID, entities.Target(note.Target), note.PublicIdentifier)
		if err != nil {
			return nil, "", err
		}
		return nil, deletedNote.ID.String(), nil
	}

	// Attempt to create a note.
	createdNote, err := s.createNoteRepository.CreateNote(
		ctx,
		note.AuthorID,
		entities.Target(note.Target),
		note.PublicIdentifier,
		&dao.CreateNoteData{
			Content:   note.Content,
			UpdatedAt: note.UpdatedAt,
		},
	)

	// Note was successfully created.
	if err == nil {
		return &models.Note{
			ID:               createdNote.ID.String(),
			PublicIdentifier: createdNote.PublicIdentifier,
			AuthorID:         createdNote.AuthorID,
			Target:           string(createdNote.Target),
			Content:          createdNote.Content,
			UpdatedAt:        createdNote.UpdatedAt,
		}, createdNote.ID.String(), nil
	}

	if !errors.Is(err, dao.ErrNoteAlreadyExists) {
		return nil, "", err
	}

	// Note already exists. Update it.
	updatedNote, err := s.updateNoteRepository.UpdateNote(
		ctx,
		note.AuthorID,
		entities.Target(note.Target),
		note.PublicIdentifier,
		&dao.UpdateNoteData{
			Content:   note.Content,
			UpdatedAt: note.UpdatedAt,
		},
	)

	if err != nil {
		return nil, "", err
	}

	return &models.Note{
		ID:               updatedNote.ID.String(),
		PublicIdentifier: updatedNote.PublicIdentifier,
		AuthorID:         updatedNote.AuthorID,
		Target:           string(updatedNote.Target),
		Content:          updatedNote.Content,
		UpdatedAt:        updatedNote.UpdatedAt,
	}, updatedNote.ID.String(), nil
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
