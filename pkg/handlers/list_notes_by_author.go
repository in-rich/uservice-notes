package handlers

import (
	"context"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListNotesByAuthorHandler struct {
	notes_pb.ListNotesByAuthorServer
	service services.ListNotesByAuthorService
}

func (h *ListNotesByAuthorHandler) ListNotesByAuthor(ctx context.Context, in *notes_pb.ListNotesByAuthorRequest) (*notes_pb.ListNotesByAuthorResponse, error) {
	notes, err := h.service.Exec(ctx, &models.ListNotesByAuthor{
		AuthorID: in.GetAuthorId(),
		Limit:    int(in.GetLimit()),
		Offset:   int(in.GetOffset()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list author notes: %v", err)
	}

	res := &notes_pb.ListNotesByAuthorResponse{
		Notes: make([]*notes_pb.Note, len(notes)),
	}
	for i, note := range notes {
		res.Notes[i] = &notes_pb.Note{
			PublicIdentifier: note.PublicIdentifier,
			AuthorId:         note.AuthorID,
			Target:           note.Target,
			Content:          note.Content,
			UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
		}
	}

	return res, nil
}

func NewListNotesByAuthorHandler(service services.ListNotesByAuthorService) *ListNotesByAuthorHandler {
	return &ListNotesByAuthorHandler{
		service: service,
	}
}
