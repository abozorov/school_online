package handler

import (
	"github.com/abozorov/school_online/cmd/classroom/internal/models"
	classroomv1 "github.com/abozorov/school_online/grpc_api/generate/classroompb/classroom/v1"
	"github.com/abozorov/school_online/pkg/logger"
)

type Handler struct {
	classroomv1.UnimplementedClassroomServiceServer
	logger  *logger.Logger
	service models.ClassroomService
}

func New(logger *logger.Logger, service models.ClassroomService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}
