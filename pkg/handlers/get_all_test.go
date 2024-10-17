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

func TestGetAllNotes(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.GetAllNotesRequest

		getAllResponse []*models.Note
		getAllErr      error

		expect     *notes_pb.GetAllNotesResponse
		expectCode codes.Code
	}{
		{
			name: "GetAllNotes",
			in: &notes_pb.GetAllNotesRequest{
				Limit:  50,
				Offset: 10,
			},
			getAllResponse: []*models.Note{
				{
					ID:               "id-1",
					PublicIdentifier: "public-identifier-1",
					AuthorID:         "author-id-1",
					Target:           "company",
					Content:          "content-1",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:               "id-2",
					PublicIdentifier: "public-identifier-2",
					AuthorID:         "author-id-1",
					Target:           "user",
					Content:          "content-2",
					UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: &notes_pb.GetAllNotesResponse{
				Notes: []*notes_pb.Note{
					{
						NoteId:           "id-1",
						PublicIdentifier: "public-identifier-1",
						AuthorId:         "author-id-1",
						Target:           "company",
						Content:          "content-1",
						UpdatedAt:        timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
					{
						NoteId:           "id-2",
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
			name: "GetAllError",
			in: &notes_pb.GetAllNotesRequest{
				Limit:  50,
				Offset: 10,
			},
			getAllErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetAllNotesService(t)

			service.
				On("Exec", context.TODO(), &models.GetAllNotes{
					Offset: tt.in.Offset,
					Limit:  tt.in.Limit,
				}).
				Return(tt.getAllResponse, tt.getAllErr)

			handler := handlers.NewGetAllNotesHandler(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.GetAllNotes(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
