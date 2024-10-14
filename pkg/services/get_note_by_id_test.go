package services_test

import (
	"context"
	"github.com/google/uuid"
	daomocks "github.com/in-rich/uservice-notes/pkg/dao/mocks"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetNoteByID(t *testing.T) {
	testData := []struct {
		name string

		selector *models.GetNoteByID

		shouldCallGetNoteByID bool
		getNoteByIDResponse   *entities.Note
		getNoteByIDError      error

		expect    *models.Note
		expectErr error
	}{
		{
			name: "ListNotes",
			selector: &models.GetNoteByID{
				AuthorID: "author-id",
				NoteID:   "00000000-0000-0000-0000-000000000001",
			},
			shouldCallGetNoteByID: true,
			getNoteByIDResponse: &entities.Note{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				PublicIdentifier: "public-identifier",
				AuthorID:         "author-id",
				Target:           entities.Target("target"),
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Note{
				ID:               "00000000-0000-0000-0000-000000000001",
				PublicIdentifier: "public-identifier",
				AuthorID:         "author-id",
				Target:           "target",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "ListNotesError",
			selector: &models.GetNoteByID{
				AuthorID: "author-id",
				NoteID:   "00000000-0000-0000-0000-000000000001",
			},
			shouldCallGetNoteByID: true,
			getNoteByIDError:      FooErr,
			expectErr:             FooErr,
		},
		{
			name: "InvalidSelector",
			selector: &models.GetNoteByID{
				AuthorID: "author-id",
				NoteID:   "",
			},
			expectErr: services.ErrInvalidNoteSelector,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getNote := daomocks.NewMockGetNoteByIDRepository(t)

			if tt.shouldCallGetNoteByID {
				getNote.
					On("GetNoteByID", context.TODO(), tt.selector.AuthorID, uuid.MustParse(tt.selector.NoteID)).
					Return(tt.getNoteByIDResponse, tt.getNoteByIDError)
			}

			service := services.NewGetNoteByIDService(getNote)

			result, err := service.Exec(context.TODO(), tt.selector)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, result)

			getNote.AssertExpectations(t)
		})
	}
}
