package handlers

import (
	"context"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListNotesByAuthorHandler struct {
	notes_pb.ListNotesByAuthorServer
	service services.ListNotesByAuthorService
	logger  monitor.GRPCLogger
}

func (h *ListNotesByAuthorHandler) listNotesByAuthor(ctx context.Context, in *notes_pb.ListNotesByAuthorRequest) (*notes_pb.ListNotesByAuthorResponse, error) {
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

func (h *ListNotesByAuthorHandler) ListNotesByAuthor(ctx context.Context, in *notes_pb.ListNotesByAuthorRequest) (*notes_pb.ListNotesByAuthorResponse, error) {
	res, err := h.listNotesByAuthor(ctx, in)
	h.logger.Report(ctx, "ListNotesByAuthor", err)
	return res, err
}

func NewListNotesByAuthorHandler(service services.ListNotesByAuthorService, logger monitor.GRPCLogger) *ListNotesByAuthorHandler {
	return &ListNotesByAuthorHandler{
		service: service,
		logger:  logger,
	}
}
