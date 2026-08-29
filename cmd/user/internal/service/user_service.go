package service

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/user/internal/models"
)

type Service struct {
	repo models.UserRepository
}

func New(repo models.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id int32) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user_service.Get: %w", models.ErrInvalidID)
	}
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user_service.Get: %w", err)
	}
	return u, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("user_service.GetByEmail: %w", models.ErrEmptyEmail)
	}
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user_service.GetByEmail: %w", err)
	}
	return u, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*models.User, error) {
	u, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("user_service.GetAll: %w", err)
	}
	return u, nil
}

func (s *Service) Create(ctx context.Context, user *models.User) (int32, error) {
	if user == nil {
		return 0, fmt.Errorf("user_service.Create: %w", models.ErrInvalidID)
	}
	err := validateUser(user, true)
	if err != nil {
		return 0, fmt.Errorf("user_service.Create: %w", err)
	}
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("user_service.Create: %w", err)
	}
	return id, nil
}

func (s *Service) Update(ctx context.Context, user *models.User) error {
	if user == nil || user.ID <= 0 {
		return fmt.Errorf("user_service.Update: %w", models.ErrInvalidID)
	}

	err := validateUser(user, false)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int32) error {
	if id <= 0 {
		return fmt.Errorf("user_service.Delete: %w", models.ErrInvalidID)
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("user_service.Delete: %w", err)
	}
	return nil
}

func validateUser(user *models.User, requirePassword bool) error {
	if user == nil {
		return models.ErrInvalidID
	}
	if user.Name == "" && requirePassword {
		return models.ErrEmptyName
	}
	if user.Username == "" && requirePassword {
		return models.ErrInvalidUsername
	}
	if user.Email == "" && requirePassword {
		return models.ErrEmptyEmail
	}
	if requirePassword && user.PasswordHash == "" {
		return models.ErrInvalidPassword
	}
	if user.Role != "" {
		switch user.Role {
		case "user":
		case "staff":
			if user.StaffRole == nil {
				return models.ErrInvalidRole
			}
		case "teacher":
			if user.TeacherRole == nil {
				return models.ErrInvalidRole
			}
		case "student":
			if user.StudentRole == nil {
				return models.ErrInvalidRole
			}
		case "parent":
			if user.ParentRole == nil {
				return models.ErrInvalidRole
			}
		default:
			return models.ErrInvalidRole
		}
	}
	if user.Birthday == "" && requirePassword {
		return models.ErrInvalidBirthday
	}
	return nil
}
