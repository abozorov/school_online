package services

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/config"
	classroom1 "github.com/abozorov/school_online/grpc_api/generate/classroompb/classroom/v1"
	raiting1 "github.com/abozorov/school_online/grpc_api/generate/raitingpb/raiting/v1"
	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
)

type IServiceManager interface {
	UserService() userv1.UserServiceClient
	ClassroomService() classroom1.ClassroomServiceClient
	RaitingService() raiting1.RaitingServiceClient
}

type serviceManager struct {
	userService      userv1.UserServiceClient
	classroomService classroom1.ClassroomServiceClient
	raitingService   raiting1.RaitingServiceClient
}

func (s *serviceManager) UserService() userv1.UserServiceClient {
	return s.userService
}

func (s *serviceManager) ClassroomService() classroom1.ClassroomServiceClient {
	return s.classroomService
}

func (s *serviceManager) RaitingService() raiting1.RaitingServiceClient {
	return s.raitingService
}

func NewServiceManager(config config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns") // dns:///localhost:50051

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.UserService.Host, config.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	connClassroomService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.ClassroomService.Host, config.ClassroomService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to classroom service: %w", err)
	}

	connRaitingService, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.RaitingService.Host, config.RaitingService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Raiting service: %w", err)
	}

	return &serviceManager{
		userService:      userv1.NewUserServiceClient(connUserService),
		classroomService: classroom1.NewClassroomServiceClient(connClassroomService),
		raitingService:   raiting1.NewRaitingServiceClient(connRaitingService),
	}, nil
}
