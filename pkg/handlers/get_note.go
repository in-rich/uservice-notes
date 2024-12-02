package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetNoteHandler struct {
	notes_pb.GetNoteServer
	service services.GetNoteService
	logger  monitor.GRPCLogger
}

func (h *GetNoteHandler) getNote(ctx context.Context, in *notes_pb.GetNoteRequest) (*notes_pb.Note, error) {
	note, err := h.service.Exec(ctx, &models.GetNote{
		Target:           in.GetTarget(),
		PublicIdentifier: in.GetPublicIdentifier(),
		AuthorID:         in.GetAuthorId(),
	})
	if err != nil {
		if errors.Is(err, dao.ErrNoteNotFound) {
			return nil, status.Error(codes.NotFound, "note not found")
		}
		if errors.Is(err, services.ErrInvalidNoteSelector) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note selector: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to get note: %v", err)
	}

	return &notes_pb.Note{
		NoteId:           note.ID,
		PublicIdentifier: note.PublicIdentifier,
		AuthorId:         note.AuthorID,
		Target:           note.Target,
		Content:          note.Content,
		UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
	}, nil
}

func (h *GetNoteHandler) GetNote(ctx context.Context, in *notes_pb.GetNoteRequest) (*notes_pb.Note, error) {
	res, err := h.getNote(ctx, in)
	h.logger.Report(ctx, "GetNote", err)
	return res, err
}

func NewGetNoteHandler(service services.GetNoteService, logger monitor.GRPCLogger) *GetNoteHandler {
	return &GetNoteHandler{
		service: service,
		logger:  logger,
	}
}
