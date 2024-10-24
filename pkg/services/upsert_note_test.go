package services_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-notes/pkg/dao"
	daomocks "github.com/in-rich/uservice-notes/pkg/dao/mocks"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpsertNote(t *testing.T) {
	testData := []struct {
		name string

		note *models.UpsertNote

		shouldCallCountUpdates bool
		countUpdatesResponse   int
		countUpdatesError      error

		shouldCallDeleteNote bool
		deleteNoteResponse   *entities.Note
		deleteNoteError      error

		shouldCallCreateNote bool
		createNoteResponse   *entities.Note
		createNoteError      error

		shouldCallUpdateNote bool
		updateNoteResponse   *entities.Note
		updateNoteError      error

		expect    *models.Note
		expectID  string
		expectErr error
	}{
		{
			name: "UpdateNote",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
			},
			shouldCallCreateNote: true,
			createNoteError:      dao.ErrNoteAlreadyExists,
			shouldCallUpdateNote: true,
			updateNoteResponse: &entities.Note{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				AuthorID:         "author-id",
				Target:           entities.TargetCompany,
				PublicIdentifier: "public-identifier",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Note{
				ID:               "00000000-0000-0000-0000-000000000001",
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectID: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "CreateNote",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
			},
			shouldCallCreateNote: true,
			createNoteResponse: &entities.Note{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				AuthorID:         "author-id",
				Target:           entities.TargetCompany,
				PublicIdentifier: "public-identifier",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Note{
				ID:               "00000000-0000-0000-0000-000000000001",
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectID: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "DeleteNote",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "",
			},
			shouldCallDeleteNote: true,
			deleteNoteResponse: &entities.Note{
				ID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},
			expectID: "00000000-0000-0000-0000-000000000001",
		},
		{
			name: "UpdateNoteError",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
			},
			shouldCallCreateNote: true,
			createNoteError:      dao.ErrNoteAlreadyExists,
			shouldCallUpdateNote: true,
			updateNoteError:      FooErr,
			expectErr:            FooErr,
		},
		{
			name: "CreateNoteError",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "content",
			},
			shouldCallCreateNote: true,
			createNoteError:      FooErr,
			expectErr:            FooErr,
		},
		{
			name: "DeleteNoteError",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "company",
				PublicIdentifier: "public-identifier",
				Content:          "",
			},
			shouldCallDeleteNote: true,
			deleteNoteError:      FooErr,
			expectErr:            FooErr,
		},
		{
			name: "InvalidTarget",
			note: &models.UpsertNote{
				AuthorID:         "author-id",
				Target:           "invalid",
				PublicIdentifier: "public-identifier",
				Content:          "content",
			},
			expectErr: services.ErrInvalidNoteUpdate,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			deleteNote := daomocks.NewMockDeleteNoteRepository(t)
			createNote := daomocks.NewMockCreateNoteRepository(t)
			updateNote := daomocks.NewMockUpdateNoteRepository(t)

			if tt.shouldCallDeleteNote {
				deleteNote.
					On("DeleteNote", context.TODO(), tt.note.AuthorID, entities.Target(tt.note.Target), tt.note.PublicIdentifier).
					Return(tt.deleteNoteResponse, tt.deleteNoteError)
			}

			if tt.shouldCallCreateNote {
				createNote.
					On(
						"CreateNote",
						context.TODO(),
						tt.note.AuthorID,
						entities.Target(tt.note.Target),
						tt.note.PublicIdentifier,
						&dao.CreateNoteData{Content: tt.note.Content},
					).
					Return(tt.createNoteResponse, tt.createNoteError)
			}

			if tt.shouldCallUpdateNote {
				updateNote.
					On(
						"UpdateNote",
						context.TODO(),
						tt.note.AuthorID,
						entities.Target(tt.note.Target),
						tt.note.PublicIdentifier,
						&dao.UpdateNoteData{Content: tt.note.Content},
					).
					Return(tt.updateNoteResponse, tt.updateNoteError)
			}

			service := services.NewUpsertNoteService(updateNote, createNote, deleteNote)

			note, noteID, err := service.Exec(context.TODO(), tt.note)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
			require.Equal(t, tt.expectID, noteID)

			deleteNote.AssertExpectations(t)
			createNote.AssertExpectations(t)
			updateNote.AssertExpectations(t)
		})
	}
}
