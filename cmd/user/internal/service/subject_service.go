package service

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/user/internal/models"
)

func (s *Service) CreateSubject(ctx context.Context, name string, description string) (int32, error) {
	if name == "" {
		return 0, fmt.Errorf("user_service.CreateSubject: %w", models.ErrInvalidName)
	}
	if description == "" {
		return 0, fmt.Errorf("user_service.CreateSubject: %w", models.ErrInvalidDescription)
	}
	id, err := s.repo.CreateSubject(ctx, name, description)
	if err != nil {
		return 0, fmt.Errorf("user_service.CreateSubject: %w", err)
	}
	return id, nil
}

func (s *Service) GetSubjectById(ctx context.Context, id int32) (*models.Subject, error) {
	if id <= 0 {
		return nil, fmt.Errorf("user_service.GetSubjectById: %w", models.ErrInvalidID)
	}
	sub, err := s.repo.GetSubjectById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user_service.GetSubjectById: %w", err)
	}
	return sub, nil
}

func (s *Service) GetAllSubjects(ctx context.Context) ([]*models.Subject, error) {
	subs, err := s.repo.GetAllSubjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("user_service.GetAllSubjects: %w", err)
	}
	return subs, nil
}

func (s *Service) UpdateSubject(ctx context.Context, id int32, name string, description string) error {
	if id <= 0 {
		return fmt.Errorf("user_service.UpdateSubject: %w", models.ErrInvalidID)
	}
	if name == "" {
		return fmt.Errorf("user_service.CreateSubject: %w", models.ErrInvalidName)
	}
	if description == "" {
		return fmt.Errorf("user_service.CreateSubject: %w", models.ErrInvalidDescription)
	}
	err := s.repo.UpdateSubject(ctx, id, name, description)
	if err != nil {
		return fmt.Errorf("user_service.UpdateSubject: %w", err)
	}
	return nil
}
