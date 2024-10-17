package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/handlers"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	servicesmocks "github.com/in-rich/uservice-notes/pkg/services/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUpsertNote(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.UpsertNoteRequest

		upsertResponse   *models.Note
		upsertIDResponse string
		upsertErr        error

		expect     *notes_pb.UpsertNoteResponse
		expectCode codes.Code
	}{
		{
			name: "UpsertNote",
			in: &notes_pb.UpsertNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Content:          "content-1",
			},
			upsertResponse: &models.Note{
				ID:               "note-id",
				PublicIdentifier: "public-identifier-1",
				AuthorID:         "author-id-1",
				Target:           "company",
				Content:          "content-1",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			upsertIDResponse: "note-id",
			expect: &notes_pb.UpsertNoteResponse{
				Id: "note-id",
				Note: &notes_pb.Note{
					NoteId:           "note-id",
					PublicIdentifier: "public-identifier-1",
					AuthorId:         "author-id-1",
					Target:           "company",
					Content:          "content-1",
					UpdatedAt:        timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "DeleteNote",
			in: &notes_pb.UpsertNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Content:          "content-1",
			},
			upsertIDResponse: "note-id",
			expect: &notes_pb.UpsertNoteResponse{
				Id: "note-id",
			},
		},
		{
			name: "InvalidArgument",
			in: &notes_pb.UpsertNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Content:          "content-1",
			},
			upsertErr:  services.ErrInvalidNoteUpdate,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",
			in: &notes_pb.UpsertNoteRequest{
				Target:           "company",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Content:          "content-1",
			},
			upsertErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpsertNoteService(t)
			service.On("Exec", context.TODO(), &models.UpsertNote{
				Target:           tt.in.GetTarget(),
				PublicIdentifier: tt.in.GetPublicIdentifier(),
				AuthorID:         tt.in.GetAuthorId(),
				Content:          tt.in.GetContent(),
			}).Return(tt.upsertResponse, tt.upsertIDResponse, tt.upsertErr)

			handler := handlers.NewUpsertNoteHandler(service, monitor.NewDummyGRPCLogger())

			resp, err := handler.UpsertNote(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
