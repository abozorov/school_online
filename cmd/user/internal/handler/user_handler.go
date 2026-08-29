package handler

import (
	"context"

	"github.com/abozorov/school_online/cmd/user/internal/models"
	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/abozorov/school_online/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	logger  *logger.Logger
	service models.UserService
}

func New(logger *logger.Logger, service models.UserService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Get(ctx context.Context, request *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: Get", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	user, err := h.service.Get(ctx, request.GetId())
	if err != nil {
		h.logger.Error("User microservice: Get", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return toProtoUser(user), nil
}

func (h *Handler) GetByEmail(ctx context.Context, request *userv1.GetByEmailRequest) (*userv1.GetUserResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: GetByEmail", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	user, err := h.service.GetByEmail(ctx, request.GetEmail())
	if err != nil {
		h.logger.Error("User microservice: GetByEmail", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return toProtoUser(user), nil
}

func (h *Handler) GetAll(ctx context.Context, request *userv1.GetAllRequest) (*userv1.GetAllUsersResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: GetAll", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	users, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("User microservice: GetAll", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	resp := &userv1.GetAllUsersResponse{Users: make([]*userv1.GetUserResponse, 0, len(users))}
	for _, user := range users {
		resp.Users = append(resp.Users, toProtoUser(user))
	}
	return resp, nil
}

func (h *Handler) Create(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: Create", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}

	id, err := h.service.Create(ctx, toModelUser(request))
	if err != nil {
		h.logger.Error("User microservice: Create", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return &userv1.CreateUserResponse{Id: id}, nil
}

func (h *Handler) Update(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: Update", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}

	if err := h.service.Update(ctx, toUpdatedModelUser(request)); err != nil {
		h.logger.Error("User microservice: Update", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return &userv1.UpdateUserResponse{}, nil
}

func (h *Handler) Delete(ctx context.Context, request *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if request == nil {
		h.logger.Error("User microservice: Delete", zap.String("error", errs.ErrBadRequest.Error()))
		return nil, responseErr(errs.ErrBadRequest)
	}
	if err := h.service.Delete(ctx, request.GetId()); err != nil {
		h.logger.Error("User microservice: Delete", zap.String("error", err.Error()))
		return nil, responseErr(err)
	}
	return &userv1.DeleteUserResponse{}, nil
}

func toProtoUser(user *models.User) *userv1.GetUserResponse {
	if user == nil {
		return &userv1.GetUserResponse{}
	}
	resp := &userv1.GetUserResponse{
		Id:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		VerifyEmail: user.VerifyEmail,
		Password:    user.PasswordHash,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		Birthday:    user.Birthday,
	}
	if user.StudentRole != nil {
		resp.Student = &userv1.Student{ClassId: user.StudentRole.ClassroomID}
	}
	if user.ParentRole != nil {
		resp.Parent = &userv1.Parent{StudentsId: user.ParentRole.StudentsID}
	}
	if user.StaffRole != nil {
		resp.Stuff = &userv1.Stuff{Position: user.StaffRole.Position, Experience: user.StaffRole.Experience}
	}
	if user.TeacherRole != nil {
		resp.Teacher = &userv1.Teacher{SubjectsId: user.TeacherRole.SubjectsID, Experience: user.TeacherRole.Experience}
	}
	return resp
}

func toModelUser(req *userv1.CreateUserRequest) *models.User {
	user := &models.User{
		Name:         req.GetName(),
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
		PhoneNumber:  req.GetPhoneNumber(),
		Role:         req.GetRole(),
		Birthday:     req.GetBirthday(),
	}
	if req.GetStudent() != nil {
		user.StudentRole = &models.StudentRole{ClassroomID: req.GetStudent().GetClassId()}
	}
	if req.GetParent() != nil {
		user.ParentRole = &models.ParentRole{StudentsID: req.GetParent().GetStudentsId()}
	}
	if req.GetStuff() != nil {
		user.StaffRole = &models.StaffRole{Position: req.GetStuff().GetPosition(), Experience: req.GetStuff().GetExperience()}
	}
	if req.GetTeacher() != nil {
		user.TeacherRole = &models.TeacherRole{SubjectsID: req.GetTeacher().GetSubjectsId(), Experience: req.GetTeacher().GetExperience()}
	}
	return user
}

func toUpdatedModelUser(req *userv1.UpdateUserRequest) *models.User {
	user := &models.User{ID: req.GetId()}
	if req.Name != nil {
		user.Name = req.GetName()
	}
	if req.Username != nil {
		user.Username = req.GetUsername()
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.GetPhoneNumber()
	}
	if req.Role != nil {
		user.Role = req.GetRole()
	}
	if req.Birthday != nil {
		user.Birthday = req.GetBirthday()
	}
	if req.GetStudent() != nil {
		user.StudentRole = &models.StudentRole{ClassroomID: req.GetStudent().GetClassId()}
	}
	if req.GetParent() != nil {
		user.ParentRole = &models.ParentRole{StudentsID: req.GetParent().GetStudentsId()}
	}
	if req.GetStuff() != nil {
		user.StaffRole = &models.StaffRole{Position: req.GetStuff().GetPosition(), Experience: req.GetStuff().GetExperience()}
	}
	if req.GetTeacher() != nil {
		user.TeacherRole = &models.TeacherRole{SubjectsID: req.GetTeacher().GetSubjectsId(), Experience: req.GetTeacher().GetExperience()}
	}
	return user
}
