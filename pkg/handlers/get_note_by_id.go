package handlers

import (
	"context"
	"errors"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetNoteByIDHandler struct {
	notes_pb.GetNoteByIDServer
	service services.GetNoteByIDService
}

func (h *GetNoteByIDHandler) GetNoteByID(ctx context.Context, in *notes_pb.GetNoteByIDRequest) (*notes_pb.Note, error) {
	note, err := h.service.Exec(ctx, &models.GetNoteByID{
		AuthorID: in.GetAuthorId(),
		NoteID:   in.GetNoteId(),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidNoteSelector) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note selector: %v", err)
		}
		if errors.Is(err, dao.ErrNoteNotFound) {
			return nil, status.Errorf(codes.NotFound, "note not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get note: %v", err)
	}

	return &notes_pb.Note{
		Target:           note.Target,
		PublicIdentifier: note.PublicIdentifier,
		Content:          note.Content,
		AuthorId:         note.AuthorID,
		NoteId:           note.ID,
		UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
	}, nil
}

func NewGetNoteByIDHandler(service services.GetNoteByIDService) *GetNoteByIDHandler {
	return &GetNoteByIDHandler{
		service: service,
	}
}
