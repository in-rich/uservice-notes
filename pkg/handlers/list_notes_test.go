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

func TestListNotes(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.ListNotesRequest

		listResponse []*models.Note
		listErr      error

		expect     *notes_pb.ListNotesResponse
		expectCode codes.Code
	}{
		{
			name: "ListNotes",
			in: &notes_pb.ListNotesRequest{
				Filters: []*notes_pb.ListNoteFilter{
					{
						Target:           "company",
						PublicIdentifier: "public-identifier-1",
					},
					{
						Target:           "user",
						PublicIdentifier: "public-identifier-2",
					},
				},
				AuthorId: "author-id-1",
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
			expect: &notes_pb.ListNotesResponse{
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
			in: &notes_pb.ListNotesRequest{
				Filters: []*notes_pb.ListNoteFilter{
					{
						Target:           "company",
						PublicIdentifier: "public-identifier-1",
					},
					{
						Target:           "user",
						PublicIdentifier: "public-identifier-2",
					},
				},
				AuthorId: "author-id-1",
			},
			listErr:    errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockListNotesService(t)

			service.
				On("Exec", context.TODO(), &models.ListNotes{
					Filters: lo.Map(tt.in.Filters, func(item *notes_pb.ListNoteFilter, index int) models.ListNotesFilter {
						return models.ListNotesFilter{
							Target:           item.GetTarget(),
							PublicIdentifier: item.GetPublicIdentifier(),
						}
					}),
					AuthorID: tt.in.AuthorId,
				}).
				Return(tt.listResponse, tt.listErr)

			handler := handlers.NewListNotesHandler(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.ListNotes(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
