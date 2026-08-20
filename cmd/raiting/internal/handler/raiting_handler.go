package handler

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/raiting/internal/models"
	raiting1 "github.com/abozorov/school_online/grpc_api/generate/raitingpb/raiting/v1"
	"github.com/abozorov/school_online/pkg/logger"
)

type Handler struct {
	raiting1.UnimplementedRaitingServiceServer
	logger  *logger.Logger
	service models.RaitingService
}

func New(logger *logger.Logger, service models.RaitingService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) GetStudentJournal(ctx context.Context, req *raiting1.GetJournalRequest) (*raiting1.GetStudentJournalResponse, error) {
	if req == nil {
		return nil, responseErr(fmt.Errorf("empty request"))
	}
	res, err := h.service.GetStudentJournal(ctx, int(req.Id), req.Period)
	if err != nil {
		return nil, responseErr(err)
	}
	return res, nil
}

func (h *Handler) GetClassroomJournal(ctx context.Context, req *raiting1.GetJournalRequest) (*raiting1.GetClassroomJournalResponse, error) {
	if req == nil {
		return nil, responseErr(fmt.Errorf("empty request"))
	}
	res, err := h.service.GetClassroomJournal(ctx, int(req.Id), req.Period)
	if err != nil {
		return nil, responseErr(err)
	}
	return res, nil
}

func (h *Handler) UpdateGrade(ctx context.Context, req *raiting1.UpdateJournalRequest) (*raiting1.UpdateJournalResponse, error) {
	if req == nil {
		return nil, responseErr(fmt.Errorf("empty request"))
	}
	res, err := h.service.UpdateGrade(ctx, req)
	if err != nil {
		return nil, responseErr(err)
	}
	return res, nil
}
