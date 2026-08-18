package userservice

import (
	"context"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
)

func (u *UserService) CreateSubject(ctx context.Context, request models.CreateSubjectRequest) (*models.Subject, error) {
	// validate request
	if err := models.ValidateCreateSubjectRequest(&request); err != nil {
		return &models.Subject{}, err
	}

	// create subject
	subject, err := u.serviceManager.UserService().CreateSubject(ctx, &userv1.CreateSubjectRequest{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return &models.Subject{}, err
	}

	return &models.Subject{
		ID:   subject.GetId(),
	}, nil
}
