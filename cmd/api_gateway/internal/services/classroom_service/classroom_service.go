package classroomservice

import "github.com/abozorov/school_online/cmd/api_gateway/internal/services"

type ClassroomService struct {
	serviceManager services.IServiceManager
}

func NewClassroomService(serviceManager services.IServiceManager) *ClassroomService {

	return &ClassroomService{
		serviceManager: serviceManager,
	}
}
