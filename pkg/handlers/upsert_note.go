package handlers

import (
	"context"
	"errors"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpsertNoteHandler struct {
	notes_pb.UpsertNoteServer
	service services.UpsertNoteService
}

func (h *UpsertNoteHandler) UpsertNote(ctx context.Context, in *notes_pb.UpsertNoteRequest) (*notes_pb.UpsertNoteResponse, error) {
	note, err := h.service.Exec(ctx, &models.UpsertNote{
		Target:           in.GetTarget(),
		PublicIdentifier: in.GetPublicIdentifier(),
		AuthorID:         in.GetAuthorId(),
		Content:          in.GetContent(),
	})
	if err != nil {
		if errors.Is(err, services.ErrNotesUpdateLimitReached) {
			return nil, status.Error(codes.ResourceExhausted, "note update limit reached")
		}
		if errors.Is(err, services.ErrInvalidNoteUpdate) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert note: %v", err)
	}

	if note == nil {
		return &notes_pb.UpsertNoteResponse{}, nil
	}

	return &notes_pb.UpsertNoteResponse{
		Note: &notes_pb.Note{
			Target:           note.Target,
			PublicIdentifier: note.PublicIdentifier,
			Content:          note.Content,
			UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
			AuthorId:         note.AuthorID,
		},
	}, nil
}

func NewUpsertNoteHandler(service services.UpsertNoteService) *UpsertNoteHandler {
	return &UpsertNoteHandler{
		service: service,
	}
}
