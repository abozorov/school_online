package handler

import (
	"context"

	"github.com/abozorov/school_online/cmd/user/internal/models"
	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/school_online/pkg/errs"
	"go.uber.org/zap"
)

func (h *Handler) CreateSubject(ctx context.Context, request *userv1.CreateSubjectRequest) (*userv1.CreateSubjectResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: CreateSubject", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	id, err := h.service.CreateSubject(ctx, request.GetName(), request.GetDescription())
	if err != nil {
		h.logger.Error("User microservice: CreateSubject", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return &userv1.CreateSubjectResponse{Id: id}, nil
}

func (h *Handler) GetSubjectById(ctx context.Context, request *userv1.GetSubjectByIdRequest) (*userv1.GetSubjectByIdResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: GetSubjectById", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	subject, err := h.service.GetSubjectById(ctx, request.GetId())
	if err != nil {
		h.logger.Error("User microservice: GetSubjectById", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return toProtoSubject(subject), nil
}

func (h *Handler) GetAllSubjects(ctx context.Context, request *userv1.GetAllRequest) (*userv1.GetAllSubjectsResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: GetAllSubjects", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	subjects, err := h.service.GetAllSubjects(ctx)
	if err != nil {
		h.logger.Error("User microservice: GetAllSubjects", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	resp := &userv1.GetAllSubjectsResponse{Subjects: make([]*userv1.GetSubjectByIdResponse, 0, len(subjects))}
	for _, subject := range subjects {
		resp.Subjects = append(resp.Subjects, toProtoSubject(subject))
	}
	return resp, nil
}

func (h *Handler) UpdateSubject(ctx context.Context, request *userv1.UpdateSubjectRequest) (*userv1.UpdateSubjectResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: UpdateSubject", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	if err := h.service.UpdateSubject(ctx, request.GetId(), request.GetName(), request.GetDescription()); err != nil {
		h.logger.Error("User microservice: UpdateSubject", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return &userv1.UpdateSubjectResponse{}, nil
}

func toProtoSubject(subject *models.Subject) *userv1.GetSubjectByIdResponse {
	if subject == nil {
		return &userv1.GetSubjectByIdResponse{}
	}
	return &userv1.GetSubjectByIdResponse{
		Id:          subject.ID,
		Name:        subject.Name,
		Description: subject.Description,
	}
}
