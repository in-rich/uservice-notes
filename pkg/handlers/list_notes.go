package handlers

import (
	"context"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListNotesHandler struct {
	notes_pb.ListNotesServer
	service services.ListNotesService
}

func (h *ListNotesHandler) ListNotes(ctx context.Context, in *notes_pb.ListNotesRequest) (*notes_pb.ListNotesResponse, error) {
	notes, err := h.service.Exec(ctx, &models.ListNotes{
		Filters: lo.Map(in.GetFilters(), func(item *notes_pb.ListNoteFilter, index int) models.ListNotesFilter {
			return models.ListNotesFilter{
				Target:           item.GetTarget(),
				PublicIdentifier: item.GetPublicIdentifier(),
			}
		}),
		AuthorID: in.GetAuthorId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	res := &notes_pb.ListNotesResponse{
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

func NewListNotesHandler(service services.ListNotesService) *ListNotesHandler {
	return &ListNotesHandler{
		service: service,
	}
}
