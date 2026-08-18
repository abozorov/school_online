package handlers

import (
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	classroomservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/classroom_service"
	raitingservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/raiting_service"
	userservice "github.com/abozorov/school_online/cmd/api_gateway/internal/services/user_service"
	"github.com/abozorov/school_online/pkg/logger"
)

type Handler struct {
	serviceManager   services.IServiceManager
	userService      *userservice.UserService
	raitingService   *raitingservice.RaitingService
	classroomService *classroomservice.ClassroomService
	logger           *logger.Logger
}

func NewHandler(serviceManager services.IServiceManager,
	userService *userservice.UserService,
	raitingService *raitingservice.RaitingService,
	classroomService *classroomservice.ClassroomService,
	logger *logger.Logger) *Handler {
	return &Handler{
		serviceManager:   serviceManager,
		userService:      userService,
		raitingService:   raitingService,
		classroomService: classroomService,
		logger:           logger,
	}
}
