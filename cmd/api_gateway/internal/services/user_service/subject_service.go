package userservice

import (
	"context"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
)

func (u *UserService) CreateSubject(ctx context.Context, request models.CreateSubjectRequest) (*models.Subject, error) {
	// validate request
	if err := models.ValidateCreateSubjectRequest(&request); err != nil {
		return &models.Subject{}, services.GRPCToErrs(err)
	}

	// create subject
	subject, err := u.serviceManager.UserService().CreateSubject(ctx, &userv1.CreateSubjectRequest{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return &models.Subject{}, services.GRPCToErrs(err)
	}

	return &models.Subject{
		ID: subject.GetId(),
	}, nil
}

func (u *UserService) GetSubjectById(ctx context.Context, id int) (*models.Subject, error) {
	if err := models.ValidateID(int32(id)); err != nil {
		return &models.Subject{}, services.GRPCToErrs(err)
	}

	subject, err := u.serviceManager.UserService().GetSubjectById(ctx, &userv1.GetSubjectByIdRequest{
		Id: int32(id),
	})
	if err != nil {
		return &models.Subject{}, services.GRPCToErrs(err)
	}

	return &models.Subject{
		ID:          subject.GetId(),
		Name:        subject.GetName(),
		Description: subject.GetDescription(),
	}, nil
}

func (u *UserService) ListSubjects(ctx context.Context) ([]*models.Subject, error) {
	subjects, err := u.serviceManager.UserService().GetAllSubjects(ctx, &userv1.GetAllRequest{})
	if err != nil {
		return nil, services.GRPCToErrs(err)
	}

	result := make([]*models.Subject, 0, len(subjects.GetSubjects()))
	for _, subject := range subjects.GetSubjects() {
		result = append(result, &models.Subject{
			ID:          subject.GetId(),
			Name:        subject.GetName(),
			Description: subject.GetDescription(),
		})
	}

	return result, nil
}

func (u *UserService) UpdateSubject(ctx context.Context, request models.UpdateSubjectRequest) error {
	if err := models.ValidateUpdateSubjectRequest(&request); err != nil {
		return services.GRPCToErrs(err)
	}

	updateRequest := &userv1.UpdateSubjectRequest{Id: request.ID}
	if request.Name != nil {
		updateRequest.Name = request.Name
	}
	if request.Description != nil {
		updateRequest.Description = request.Description
	}

	_, err := u.serviceManager.UserService().UpdateSubject(ctx, updateRequest)
	if err != nil {
		return services.GRPCToErrs(err)
	}

	return nil
}
