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

type GetAllNotesHandler struct {
	notes_pb.GetAllNotesServer
	service services.GetAllNotesService
	logger  monitor.GRPCLogger
}

func (h *GetAllNotesHandler) getAllNotes(ctx context.Context, in *notes_pb.GetAllNotesRequest) (*notes_pb.GetAllNotesResponse, error) {
	notes, err := h.service.Exec(ctx, &models.GetAllNotes{Limit: in.GetLimit(), Offset: in.GetOffset()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all notes: %v", err)
	}

	res := &notes_pb.GetAllNotesResponse{
		Notes: make([]*notes_pb.Note, len(notes)),
	}
	for i, note := range notes {
		res.Notes[i] = &notes_pb.Note{
			NoteId:           note.ID,
			PublicIdentifier: note.PublicIdentifier,
			AuthorId:         note.AuthorID,
			Target:           note.Target,
			Content:          note.Content,
			UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
		}
	}

	return res, nil
}

func (h *GetAllNotesHandler) GetAllNotes(ctx context.Context, in *notes_pb.GetAllNotesRequest) (*notes_pb.GetAllNotesResponse, error) {
	res, err := h.getAllNotes(ctx, in)
	h.logger.Report(ctx, "GetAllNotes", err)
	return res, err
}

func NewGetAllNotesHandler(service services.GetAllNotesService, logger monitor.GRPCLogger) *GetAllNotesHandler {
	return &GetAllNotesHandler{
		service: service,
		logger:  logger,
	}
}
