package service

import "github.com/abozorov/school_online/cmd/classroom/internal/models"

type Service struct {
	repo models.ClassroomRepository
}

func New(repo models.ClassroomRepository) *Service {
	return &Service{
		repo: repo,
	}
}
