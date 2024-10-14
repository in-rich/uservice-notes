package handlers_test

import (
	"context"
	"errors"
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

func TestGetNoteByID(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.GetNoteByIDRequest

		getResponse *models.Note
		getErr      error

		expect     *notes_pb.Note
		expectCode codes.Code
	}{
		{
			name: "GetNoteByID",
			in: &notes_pb.GetNoteByIDRequest{
				AuthorId: "author-id-1",
				NoteId:   "note_id-1",
			},
			getResponse: &models.Note{
				ID:               "note_id-1",
				PublicIdentifier: "public-identifier-1",
				AuthorID:         "author-id-1",
				Target:           "company",
				Content:          "content-1",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &notes_pb.Note{
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Target:           "company",
				Content:          "content-1",
				UpdatedAt:        timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				NoteId:           "note_id-1",
			},
		},
		{
			name: "ListError",
			in: &notes_pb.GetNoteByIDRequest{
				AuthorId: "author-id-1",
				NoteId:   "note_id-1",
			},
			getErr:     errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetNoteByIDService(t)

			service.
				On("Exec", context.TODO(), &models.GetNoteByID{
					AuthorID: tt.in.AuthorId,
					NoteID:   tt.in.NoteId,
				}).
				Return(tt.getResponse, tt.getErr)

			handler := handlers.NewGetNoteByIDHandler(service)
			resp, err := handler.GetNoteByID(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
