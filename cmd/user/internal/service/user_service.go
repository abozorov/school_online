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
		return nil, models.ErrInvalidID
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, models.ErrEmptyEmail
	}
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetAll(ctx context.Context) ([]*models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) Create(ctx context.Context, user *models.User) (int32, error) {
	if user == nil {
		return 0, models.ErrInvalidID
	}
	if err := validateUser(user, true); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, user)
}

func (s *Service) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return models.ErrInvalidID
	}
	if user.ID <= 0 {
		return models.ErrInvalidID
	}
	if err := validateUser(user, false); err != nil {
		return err
	}
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id int32) error {
	if id <= 0 {
		return models.ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) CreateSubject(ctx context.Context, name string, description string) (int32, error) {
	if name == "" {
		return 0, models.ErrInvalidName
	}
	if description == "" {
		return 0, fmt.Errorf("empty description")
	}
	return s.repo.CreateSubject(ctx, name, description)
}

func (s *Service) UpdateSubject(ctx context.Context, id int32, name string, description string) error {
	if id <= 0 {
		return models.ErrInvalidID
	}
	return s.repo.UpdateSubject(ctx, id, name, description)
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
		case "user", "staff", "teacher", "student", "parent":
		default:
			return models.ErrInvalidRole
		}
	}
	if user.Birthday == "" && requirePassword {
		return models.ErrInvalidBirthday
	}
	return nil
}
