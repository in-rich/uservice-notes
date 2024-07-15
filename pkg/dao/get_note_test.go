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

var getNoteFixtures = []*entities.Note{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestGetNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		authorID         string
		publicIdentifier string
		target           entities.Target
		expect           *entities.Note
		expectErr        error
	}{
		{
			name:             "Exec",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetUser,
			expect: &entities.Note{
				ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				AuthorID:         "author-id-1",
				PublicIdentifier: "public-identifier-1",
				Target:           entities.TargetUser,
				Content:          "content-1",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name:             "NoteNotFound",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-2",
			target:           entities.TargetUser,
			expectErr:        dao.ErrNoteNotFound,
		},
	}

	stx := BeginTX(db, getNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewGetNoteRepository(tx)
			note, err := repo.GetNote(context.TODO(), tt.authorID, tt.target, tt.publicIdentifier)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
