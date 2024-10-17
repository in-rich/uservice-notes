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

var getAllNotesFixtures = []*entities.Note{
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
		AuthorID:         "author-id-2",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetUser,
		Content:          "content-2",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		AuthorID:         "author-id-3",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetCompany,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
		AuthorID:         "author-id-4",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetCompany,
		Content:          "content-2",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)),
	},
}

func TestGetAllNotes(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		limit  int64
		offset int64

		expect    []*entities.Note
		expectErr error
	}{
		{
			name:  "GetAllNotes",
			limit: 10,
			expect: []*entities.Note{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
					AuthorID:         "author-id-4",
					PublicIdentifier: "public-identifier-2",
					Target:           entities.TargetCompany,
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					AuthorID:         "author-id-3",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetCompany,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					AuthorID:         "author-id-2",
					PublicIdentifier: "public-identifier-2",
					Target:           entities.TargetUser,
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetUser,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name:   "GetAllNotes/NoResult",
			limit:  10,
			offset: 10,
			expect: []*entities.Note{},
		},
	}

	stx := BeginTX(db, getAllNotesFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewGetAllNotesRepository(tx)
			note, err := repo.GetAllNotes(context.TODO(), tt.limit, tt.offset)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
