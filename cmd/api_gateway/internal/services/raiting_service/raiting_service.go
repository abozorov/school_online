package raitingservice

import "github.com/abozorov/school_online/cmd/api_gateway/internal/services"

type RaitingService struct {
	serviceManager services.IServiceManager
}

func NewRaitingService(serviceManager services.IServiceManager) *RaitingService {

	return &RaitingService{
		serviceManager: serviceManager,
	}
}
