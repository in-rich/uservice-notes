package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var createNoteFixtures = []*entities.Note{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestCreateNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		authorID         string
		publicIdentifier string
		target           entities.Target
		data             *dao.CreateNoteData
		expect           *entities.Note
		expectErr        error
	}{
		{
			name:             "CreateNote/SameTarget/DifferentIdentifier/SameAuthor",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-2",
			target:           entities.TargetUser,
			data:             &dao.CreateNoteData{Content: "new-content"},
			expect: &entities.Note{
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-2",
				Target:           entities.TargetUser,
				Content:          "new-content",
			},
		},
		{
			name:             "CreateNote/DifferentTarget/SameIdentifier/SameAuthor",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetCompany,
			data:             &dao.CreateNoteData{Content: "new-content"},
			expect: &entities.Note{
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetCompany,
				Content:          "new-content",
			},
		},
		{
			name:             "CreateNote/SameTarget/SameIdentifier/DifferentAuthor",
			authorID:         "author-id-2",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetUser,
			data:             &dao.CreateNoteData{Content: "new-content"},
			expect: &entities.Note{
				AuthorID:         "author-id-2",
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetUser,
				Content:          "new-content",
			},
		},
		{
			name:             "CreateNote/SameTarget/SameIdentifier/SameAuthor",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetUser,
			data:             &dao.CreateNoteData{Content: "new-content"},
			expectErr:        dao.ErrNoteAlreadyExists,
		},
	}

	stx := BeginTX(db, createNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateNoteRepository(tx)
			note, err := repo.CreateNote(context.TODO(), tt.authorID, tt.target, tt.publicIdentifier, tt.data)

			if note != nil {
				// Since ID and UpdatedAt are random, nullify them for comparison.
				note.ID = nil
				note.UpdatedAt = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
