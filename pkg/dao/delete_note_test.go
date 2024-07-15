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

var deleteNoteFixtures = []*entities.Note{
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestDeleteNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name             string
		authorID         string
		publicIdentifier string
		target           entities.Target
		expectErr        error
	}{
		{
			name:             "DeleteNote",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-1",
			target:           entities.TargetUser,
		},
		{
			// Still a success because this method if forgiving.
			name:             "NoteNotFound",
			authorID:         "author-id-1",
			publicIdentifier: "public-identifier-2",
			target:           entities.TargetUser,
		},
	}

	stx := BeginTX(db, deleteNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteNoteRepository(tx)
			err := repo.DeleteNote(context.TODO(), tt.authorID, tt.target, tt.publicIdentifier)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
