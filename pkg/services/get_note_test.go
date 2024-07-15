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

func TestGetNoteService(t *testing.T) {
	testData := []struct {
		name string

		selector *models.GetNote

		shouldCallGetNote bool
		getNoteResponse   *entities.Note
		getNoteError      error

		expect    *models.Note
		expectErr error
	}{
		{
			name: "GetNote",
			selector: &models.GetNote{
				AuthorID:         "author-id",
				Target:           "user",
				PublicIdentifier: "public-identifier",
			},
			shouldCallGetNote: true,
			getNoteResponse: &entities.Note{
				PublicIdentifier: "public-identifier",
				AuthorID:         "author-id",
				Target:           entities.Target("target"),
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Note{
				PublicIdentifier: "public-identifier",
				AuthorID:         "author-id",
				Target:           "target",
				Content:          "content",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "GetNoteError",
			selector: &models.GetNote{
				AuthorID:         "author-id",
				Target:           "user",
				PublicIdentifier: "public-identifier",
			},
			shouldCallGetNote: true,
			getNoteError:      FooErr,
			expectErr:         FooErr,
		},
		{
			name: "InvalidSelector",
			selector: &models.GetNote{
				AuthorID:         "author-id",
				Target:           "invalid",
				PublicIdentifier: "public-identifier",
			},
			expectErr: services.ErrInvalidNoteSelector,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			getNote := daomocks.NewMockGetNoteRepository(t)

			if tt.shouldCallGetNote {
				getNote.
					On("GetNote", context.TODO(), tt.selector.AuthorID, entities.Target(tt.selector.Target), tt.selector.PublicIdentifier).
					Return(tt.getNoteResponse, tt.getNoteError)
			}

			service := services.NewGetNoteService(getNote)

			note, err := service.Exec(context.TODO(), tt.selector)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)

			getNote.AssertExpectations(t)
		})
	}
}
