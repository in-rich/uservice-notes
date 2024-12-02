package handlers

import (
	"context"
	"errors"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/pkg/models"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UpsertNoteHandler struct {
	notes_pb.UpsertNoteServer
	service services.UpsertNoteService
	logger  monitor.GRPCLogger
}

func (h *UpsertNoteHandler) upsertNote(ctx context.Context, in *notes_pb.UpsertNoteRequest) (*notes_pb.UpsertNoteResponse, error) {
	note, noteID, err := h.service.Exec(ctx, &models.UpsertNote{
		Target:           in.GetTarget(),
		PublicIdentifier: in.GetPublicIdentifier(),
		AuthorID:         in.GetAuthorId(),
		Content:          in.GetContent(),
		UpdatedAt: lo.TernaryF[*time.Time](
			in.GetUpdatedAt() == nil,
			func() *time.Time {
				return nil
			},
			func() *time.Time {
				return lo.ToPtr(in.GetUpdatedAt().AsTime())
			},
		),
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
		return &notes_pb.UpsertNoteResponse{
			Id: noteID,
		}, nil
	}

	return &notes_pb.UpsertNoteResponse{
		Id: noteID,
		Note: &notes_pb.Note{
			NoteId:           note.ID,
			Target:           note.Target,
			PublicIdentifier: note.PublicIdentifier,
			Content:          note.Content,
			UpdatedAt:        TimeToTimestampProto(note.UpdatedAt),
			AuthorId:         note.AuthorID,
		},
	}, nil
}

func (h *UpsertNoteHandler) UpsertNote(ctx context.Context, in *notes_pb.UpsertNoteRequest) (*notes_pb.UpsertNoteResponse, error) {
	res, err := h.upsertNote(ctx, in)
	h.logger.Report(ctx, "UpsertNote", err)
	return res, err
}

func NewUpsertNoteHandler(service services.UpsertNoteService, logger monitor.GRPCLogger) *UpsertNoteHandler {
	return &UpsertNoteHandler{
		service: service,
		logger:  logger,
	}
}
