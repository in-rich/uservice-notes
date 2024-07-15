package services_test

import (
	"context"
	daomocks "github.com/in-rich/uservice-notes/pkg/dao/mocks"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestListNotesByAuthor(t *testing.T) {
	testData := []struct {
		name string

		selector *models.ListNotesByAuthor

		shouldCallListNotesByAuthor bool
		listNotesByAuthorResponse   []*entities.Note
		listNotesByAuthorError      error

		expect    []*models.Note
		expectErr error
	}{
		{
			name: "ListNotes",
			selector: &models.ListNotesByAuthor{
				AuthorID: "author-id",
				Limit:    10,
			},
			shouldCallListNotesByAuthor: true,
			listNotesByAuthorResponse: []*entities.Note{
				{
					PublicIdentifier: "public-identifier",
					AuthorID:         "author-id",
					Target:           entities.Target("target"),
					Content:          "content",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: []*models.Note{
				{
					PublicIdentifier: "public-identifier",
					AuthorID:         "author-id",
					Target:           "target",
					Content:          "content",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "ListNotesError",
			selector: &models.ListNotesByAuthor{
				AuthorID: "author-id",
				Limit:    10,
			},
			shouldCallListNotesByAuthor: true,
			listNotesByAuthorError:      FooErr,
			expectErr:                   FooErr,
		},
		{
			name: "InvalidSelector",
			selector: &models.ListNotesByAuthor{
				AuthorID: "author-id",
				Limit:    -2,
			},
			expectErr: services.ErrInvalidNoteSelector,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			listNotes := daomocks.NewMockListNotesByAuthorRepository(t)

			if tt.shouldCallListNotesByAuthor {
				listNotes.
					On("ListNotesByAuthor", context.TODO(), tt.selector.AuthorID, tt.selector.Limit, tt.selector.Offset).
					Return(tt.listNotesByAuthorResponse, tt.listNotesByAuthorError)
			}

			service := services.NewListNotesByAuthorService(listNotes)

			result, err := service.Exec(context.TODO(), tt.selector)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, result)

			listNotes.AssertExpectations(t)
		})
	}
}
