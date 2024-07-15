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

var updateNoteFixtures = []*entities.Note{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestUpdateNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		authorID         string
		publicIdentifier string
		target           entities.Target
		data             *dao.UpdateNoteData
		expect           *entities.Note
		expectErr        error
	}{
		{
			name:             "UpdateNote",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetUser,
			data: &dao.UpdateNoteData{
				Content: "new-content",
			},
			expect: &entities.Note{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetUser,
				Content:          "new-content",
			},
		},
		{
			name:             "NoteNotFound",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-2",
			target:           entities.TargetUser,
			data: &dao.UpdateNoteData{
				Content: "new-content",
			},
			expectErr: dao.ErrNoteNotFound,
		},
	}

	stx := BeginTX(db, updateNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewUpdateNoteRepository(tx)
			note, err := repo.UpdateNote(context.TODO(), tt.authorID, tt.target, tt.publicIdentifier, tt.data)

			if note != nil {
				// Since UpdatedAt is random, nullify it for comparison.
				note.UpdatedAt = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
