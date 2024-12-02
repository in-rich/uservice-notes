package handlers_test

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/dao"
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

func TestGetUserData(t *testing.T) {
	testData := []struct {
		name string

		in *notes_pb.GetNoteRequest

		getResponse *models.Note
		getErr      error

		expect     *notes_pb.Note
		expectCode codes.Code
	}{
		{
			name: "GetNote",
			in: &notes_pb.GetNoteRequest{
				Target:           "user",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			getResponse: &models.Note{
				ID:               "id-1",
				PublicIdentifier: "public-identifier-1",
				AuthorID:         "author-id-1",
				Target:           "user",
				Content:          "content-1",
				UpdatedAt:        lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &notes_pb.Note{
				NoteId:           "id-1",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
				Target:           "user",
				Content:          "content-1",
				UpdatedAt:        timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "NoteNotFound",
			in: &notes_pb.GetNoteRequest{
				Target:           "user",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			getErr:     dao.ErrNoteNotFound,
			expectCode: codes.NotFound,
		},
		{
			name: "InvalidArgument",
			in: &notes_pb.GetNoteRequest{
				Target:           "user",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			getErr:     services.ErrInvalidNoteSelector,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &notes_pb.GetNoteRequest{
				Target:           "user",
				PublicIdentifier: "public-identifier-1",
				AuthorId:         "author-id-1",
			},
			getErr:     errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetNoteService(t)
			service.
				On("Exec", context.TODO(), &models.GetNote{
					Target:           tt.in.Target,
					PublicIdentifier: tt.in.PublicIdentifier,
					AuthorID:         tt.in.AuthorId,
				}).
				Return(tt.getResponse, tt.getErr)

			handler := handlers.NewGetNoteHandler(service, monitor.NewDummyGRPCLogger())
			resp, err := handler.GetNote(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
