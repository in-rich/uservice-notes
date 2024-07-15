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

var listNotesByAuthorFixtures = []*entities.Note{
	// User 1
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetUser,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-2",
		Target:           entities.TargetUser,
		Content:          "content-2",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)),
	},
	{
		ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		AuthorID:         "author-id-1",
		PublicIdentifier: "public-identifier-1",
		Target:           entities.TargetCompany,
		Content:          "content-1",
		UpdatedAt:        lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
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

func TestListNotesByAuthor(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID string
		limit    int
		offset   int

		expect    []*entities.Note
		expectErr error
	}{
		{
			name:     "ListNotesByAuthor",
			authorID: "author-id-1",
			limit:    100,
			expect: []*entities.Note{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-2",
					Target:           entities.TargetUser,
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetUser,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetCompany,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
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
			name:     "ListNotesByAuthor/LimitAndOffset",
			authorID: "author-id-1",
			limit:    2,
			offset:   1,
			expect: []*entities.Note{
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetUser,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					AuthorID:         "author-id-1",
					PublicIdentifier: "public-identifier-1",
					Target:           entities.TargetCompany,
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name:     "ListNotesByAuthor/NoResult",
			authorID: "author-id-3",
			limit:    100,
			expect:   []*entities.Note{},
		},
	}

	stx := BeginTX(db, listNotesByAuthorFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewListNotesByAuthorRepository(tx)
			note, err := repo.ListNotesByAuthor(context.TODO(), tt.authorID, tt.limit, tt.offset)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
