package userservice

import (
	"context"
	"fmt"
	"log"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	"github.com/abozorov/school_online/pkg/cache"
	"github.com/abozorov/school_online/pkg/errs"
	"github.com/abozorov/school_online/pkg/jwt"
	mailsender "github.com/abozorov/school_online/pkg/mail_sender"
	"github.com/abozorov/school_online/pkg/password"

	userv1 "github.com/abozorov/school_online/grpc_api/generate/userpb/user/v1"
)

type UserService struct {
	serviceManager services.IServiceManager
	jwt            *jwt.JWTSecret
	memCache       cache.ICache
	mailSender     *mailsender.MailSender
}

func NewUserService(
	serviceManager services.IServiceManager,
	jwt *jwt.JWTSecret,
	memCache cache.ICache,
	mailsender *mailsender.MailSender) *UserService {

	return &UserService{
		serviceManager: serviceManager,
		jwt:            jwt,
		memCache:       memCache,
		mailSender:     mailsender,
	}
}

func (u *UserService) Login(ctx context.Context, request models.LoginRequest) (*models.Tokens, error) {
	err := request.Validate()
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w: %w", errs.ErrBadRequestBody, err)
	}

	// get user by email
	user, err := u.serviceManager.UserService().GetByEmail(ctx, &userv1.GetByEmailRequest{
		Email: request.Email,
	})
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// compare password
	err = password.Compare(user.GetPassword(), request.Password)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// generate tokens
	jwtToken, err := u.jwt.GenerateToken(int(user.GetId()), user.GetEmail(), user.GetRole())
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// return tokens
	return &models.Tokens{
		JWT: jwtToken,
	}, nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := u.serviceManager.UserService().Get(ctx, &userv1.GetUserRequest{
		Id: int32(id),
	})
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", err)
	}

	return &models.User{
		ID:          user.GetId(),
		Name:        user.GetName(),
		Username:    user.GetUsername(),
		Email:       user.GetEmail(),
		VerifyEmail: user.GetVerifyEmail(),
		PhoneNumber: user.GetPhoneNumber(),
		Role:        user.GetRole(),
		Birthday:    user.GetBirthday(),
		StudentRole: models.StudentRole{
			ClassroomID: user.GetStudent().GetClassId(),
		},
		ParentRole: models.ParentRole{
			StudentsID: user.GetParent().GetStudentsId(),
		},
		StaffRole: models.StaffRole{
			Position:   user.GetStuff().GetPosition(),
			Experience: user.GetStuff().GetExperience(),
		},
		TeacherRole: models.TeacherRole{
			SubjectsID: user.GetTeacher().GetSubjectsId(),
			Experience: user.GetTeacher().GetExperience(),
		},
	}, nil
}

func (u *UserService) List(ctx context.Context) ([]*models.User, error) {
	users, err := u.serviceManager.UserService().GetAll(ctx, &userv1.GetAllUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("user_service.List: %w", err)
	}

	var result []*models.User
	for _, user := range users.GetUsers() {
		result = append(result, &models.User{
			ID:          user.GetId(),
			Name:        user.GetName(),
			Username:    user.GetUsername(),
			Email:       user.GetEmail(),
			VerifyEmail: user.GetVerifyEmail(),
			PhoneNumber: user.GetPhoneNumber(),
			Role:        user.GetRole(),
			Birthday:    user.GetBirthday(),
			StudentRole: models.StudentRole{
				ClassroomID: user.GetStudent().GetClassId(),
			},
			ParentRole: models.ParentRole{
				StudentsID: user.GetParent().GetStudentsId(),
			},
			StaffRole: models.StaffRole{
				Position:   user.GetStuff().GetPosition(),
				Experience: user.GetStuff().GetExperience(),
			},
			TeacherRole: models.TeacherRole{
				SubjectsID: user.GetTeacher().GetSubjectsId(),
				Experience: user.GetTeacher().GetExperience(),
			},
		})
	}

	return result, nil
}

func (u *UserService) Create(ctx context.Context, request models.RegisterUserRequest) (int32, error) {
	err := models.ValidateRegisterRequest(request)
	if err != nil {
		return 0, fmt.Errorf("user_service.Create: %w: %w", errs.ErrBadRequestBody, err)
	}

	// hash password
	hashedPassword, err := password.Hash(request.Password)
	if err != nil {
		return 0, fmt.Errorf("user_service.Create: %w", err)
	}

	// create user
	user := userv1.CreateUserRequest{
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    hashedPassword,
		Role:        request.Role,
		PhoneNumber: request.PhoneNumber,
		Birthday:    request.Birthday,
	}

	log.Println(request)

	switch request.Role {
	case "student":
		user.Student = &userv1.Student{
			ClassId: request.StudentRole.ClassroomID,
		}
	case "parent":
		user.Parent = &userv1.Parent{
			StudentsId: request.ParentRole.StudentsID,
		}
	case "staff":
		user.Stuff = &userv1.Stuff{
			Position:   request.StaffRole.Position,
			Experience: request.StaffRole.Experience,
		}
	case "teacher":
		user.Teacher = &userv1.Teacher{
			SubjectsId: request.TeacherRole.SubjectsID,
			Experience: request.TeacherRole.Experience,
		}
	default:
		return 0, fmt.Errorf("user_service.Create: %w", errs.ErrBadRequestBody)
	}
	out, err := u.serviceManager.UserService().Create(ctx, &user)

	if err != nil {
		return 0, fmt.Errorf("user_service.Create: %w", err)
	}

	return out.GetId(), nil
}

func (u *UserService) UpdateByID(ctx context.Context, request models.UpdateUserRequest) error {
	err := models.ValidateID(request.ID)
	if err != nil {
		return fmt.Errorf("user_service.UpdateByID: %w: %w", errs.ErrBadRequestBody, err)
	}

	// update user
	updatedUser := &userv1.UpdateUserRequest{
		Id:          request.ID,
		Name:        request.Name,
		Username:    request.Username,
		PhoneNumber: request.PhoneNumber,
		Role:        request.Role,
		Birthday:    request.Birthday,
	}

	switch *request.Role {
	case "student":
		updatedUser.Student = &userv1.Student{
			ClassId: request.StudentRole.ClassroomID,
		}
	case "parent":
		updatedUser.Parent = &userv1.Parent{
			StudentsId: request.ParentRole.StudentsID,
		}
	case "staff":
		updatedUser.Stuff = &userv1.Stuff{
			Position:   request.StaffRole.Position,
			Experience: request.StaffRole.Experience,
		}
	case "teacher":
		updatedUser.Teacher = &userv1.Teacher{
			SubjectsId: request.TeacherRole.SubjectsID,
			Experience: request.TeacherRole.Experience,
		}
	}

	_, err = u.serviceManager.UserService().Update(ctx, updatedUser)

	if err != nil {
		return fmt.Errorf("user_service.UpdateByID: %w", err)
	}

	return nil
}

func (u *UserService) DeleteByID(ctx context.Context, id int) error {
	err := models.ValidateID(int32(id))
	if err != nil {
		return fmt.Errorf("user_service.DeleteByID: %w: %w", errs.ErrBadRequestBody, err)
	}

	_, err = u.serviceManager.UserService().Delete(ctx, &userv1.DeleteUserRequest{
		Id: int32(id),
	})
	if err != nil {
		return fmt.Errorf("user_service.DeleteByID: %w", err)
	}

	return nil
}
