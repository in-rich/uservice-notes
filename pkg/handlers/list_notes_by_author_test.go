package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/handlers"
	"github.com/in-rich/uservice-notes/pkg/models"
	servicesmocks "github.com/in-rich/uservice-notes/pkg/services/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestListNotesByAuthor(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.ListNotesByAuthorRequest

		listResponse []*models.Note
		listErr      error

		expect     *notes_pb.ListNotesByAuthorResponse
		expectCode codes.Code
	}{
		{
			name: "ListNotesByAuthor",
			in: &notes_pb.ListNotesByAuthorRequest{
				AuthorId: "author-id-1",
				Limit:    50,
				Offset:   10,
			},
			listResponse: []*models.Note{
				{
					PublicIdentifier: "public-identifier-1",
					AuthorID:         "author-id-1",
					Target:           "company",
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					PublicIdentifier: "public-identifier-2",
					AuthorID:         "author-id-1",
					Target:           "user",
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: &notes_pb.ListNotesByAuthorResponse{
				Notes: []*notes_pb.Note{
					{
						PublicIdentifier: "public-identifier-1",
						AuthorId:         "author-id-1",
						Target:           "company",
						Content:          "content-1",
						UpdatedAt:        timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
					{
						PublicIdentifier: "public-identifier-2",
						AuthorId:         "author-id-1",
						Target:           "user",
						Content:          "content-2",
						UpdatedAt:        timestamppb.New(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
		},
		{
			name: "ListError",
			in: &notes_pb.ListNotesByAuthorRequest{
				AuthorId: "author-id-1",
				Limit:    50,
				Offset:   10,
			},
			listErr:    errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockListNotesByAuthorService(t)

			service.
				On("Exec", context.TODO(), &models.ListNotesByAuthor{
					Offset:   int(tt.in.Offset),
					Limit:    int(tt.in.Limit),
					AuthorID: tt.in.AuthorId,
				}).
				Return(tt.listResponse, tt.listErr)

			handler := handlers.NewListNotesByAuthorHandler(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.ListNotesByAuthor(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
