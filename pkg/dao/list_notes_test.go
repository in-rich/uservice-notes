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

var listNotesFixtures = []*entities.Note{
	// User 1
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetUser,
		Content:          "content-2",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetCompany,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetCompany,
		Content:          "content-2",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},

	// User 2
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000005")),
		AuthorID:         "author-id-2",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestListNotes(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID string
		filters  []dao.ListNoteFilter

		expect    []*entities.Note
		expectErr error
	}{
		{
			name:     "ListNotes",
			authorID: "author-id-1",
			filters: []dao.ListNoteFilter{
				{PublicIdentifier: "public-identifier-1", Target: entities.TargetUser},
				{PublicIdentifier: "public-identifier-2", Target: entities.TargetCompany},
				{PublicIdentifier: "public-identifier-3", Target: entities.TargetUser},
			},
			expect: []*entities.Note{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetUser,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-2",
					Target:           entities.TargetCompany,
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name:     "ListNotes/NoResult",
			authorID: "author-id-1",
			filters: []dao.ListNoteFilter{
				{PublicIdentifier: "public-identifier-3", Target: entities.TargetUser},
			},
			expect: []*entities.Note{},
		},
	}

	stx := BeginTX(db, listNotesFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewListNotesRepository(tx)
			note, err := repo.ListNotes(context.TODO(), tt.authorID, tt.filters)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
