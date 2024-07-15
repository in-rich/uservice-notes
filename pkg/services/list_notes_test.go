package services_test

import (
	"context"
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

func TestListNotes(t *testing.T) {
	testData := []struct {
		name string

		selector *models.ListNotes

		shouldCallListNotes bool
		listNotesResponse   []*entities.Note
		listNotesError      error

		expect    []*models.Note
		expectErr error
	}{
		{
			name: "ListNotes",
			selector: &models.ListNotes{
				AuthorID: "author-id",
				Filters: []models.ListNotesFilter{
					{
						PublicIdentifier: "public-identifier",
						Target:           "user",
					},
					{
						PublicIdentifier: "public-identifier",
						Target:           "company",
					},
				},
			},
			shouldCallListNotes: true,
			listNotesResponse: []*entities.Note{
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
			selector: &models.ListNotes{
				AuthorID: "author-id",
				Filters: []models.ListNotesFilter{
					{
						PublicIdentifier: "public-identifier",
						Target:           "user",
					},
					{
						PublicIdentifier: "public-identifier",
						Target:           "company",
					},
				},
			},
			shouldCallListNotes: true,
			listNotesError:      FooErr,
			expectErr:           FooErr,
		},
		{
			name: "InvalidSelector",
			selector: &models.ListNotes{
				AuthorID: "author-id",
				Filters: []models.ListNotesFilter{
					{
						PublicIdentifier: "public-identifier",
						Target:           "user",
					},
					{
						PublicIdentifier: "public-identifier",
						Target:           "invalid",
					},
				},
			},
			expectErr: services.ErrInvalidNoteSelector,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			listNotes := daomocks.NewMockListNotesRepository(t)

			if tt.shouldCallListNotes {
				listNotes.
					On(
						"ListNotes",
						context.TODO(),
						tt.selector.AuthorID,
						lo.Map(tt.selector.Filters, func(item models.ListNotesFilter, index int) dao.ListNoteFilter {
							return dao.ListNoteFilter{
								PublicIdentifier: item.PublicIdentifier,
								Target:           entities.Target(item.Target),
							}
						}),
					).
					Return(tt.listNotesResponse, tt.listNotesError)
			}

			service := services.NewListNotesService(listNotes)

			result, err := service.Exec(context.TODO(), tt.selector)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, result)

			listNotes.AssertExpectations(t)
		})
	}
}
