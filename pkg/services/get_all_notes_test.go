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

func TestGetAllNotes(t *testing.T) {
	testData := []struct {
		name string

		selector *models.GetAllNotes

		shouldCallGetAllNotes bool
		getAllNotesResponse   []*entities.Note
		getAllNotesError      error

		expect    []*models.Note
		expectErr error
	}{
		{
			name: "GetAllNotes",
			selector: &models.GetAllNotes{
				Limit: 10,
			},
			shouldCallGetAllNotes: true,
			getAllNotesResponse: []*entities.Note{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					PublicIdentifier: "public-identifier",
					AuthorID:         "author-id",
					Target:           entities.Target("target"),
					Content:          "content",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: []*models.Note{
				{
					ID:               "00000000-0000-0000-0000-000000000001",
					PublicIdentifier: "public-identifier",
					AuthorID:         "author-id",
					Target:           "target",
					Content:          "content",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "GetAllNotesError",
			selector: &models.GetAllNotes{
				Limit: 10,
			},
			shouldCallGetAllNotes: true,
			getAllNotesError:      FooErr,
			expectErr:             FooErr,
		},
		{
			name: "InvalidSelector",
			selector: &models.GetAllNotes{
				Limit: -2,
			},
			expectErr: services.ErrInvalidNoteSelector,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getAllNotes := daomocks.NewMockGetAllNotesRepository(t)

			if tt.shouldCallGetAllNotes {
				getAllNotes.
					On("GetAllNotes", context.TODO(), tt.selector.Limit, tt.selector.Offset).
					Return(tt.getAllNotesResponse, tt.getAllNotesError)
			}

			service := services.NewGetAllNotesService(getAllNotes)

			result, err := service.Exec(context.TODO(), tt.selector)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, result)

			getAllNotes.AssertExpectations(t)
		})
	}
}
