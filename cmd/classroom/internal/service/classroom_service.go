package service

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/classroom/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
)

type Service struct {
	repo models.ClassroomRepository
}

func New(repo models.ClassroomRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetClassroom(ctx context.Context, id int32) (*models.Classroom, error) {
	if err := models.ValidateClassroomID(id); err != nil {
		return nil, fmt.Errorf("classroom_service.GetClassroom: %w: %w", errs.ErrBadRequest, err)
	}
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetClassroom: %w", err)
	}
	return c, nil
}

func (s *Service) ListClassrooms(ctx context.Context, page, limit int32) ([]*models.Classroom, error) {
	list, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.ListClassrooms: %w", err)
	}
	return list, nil
}

func (s *Service) CreateClassroom(ctx context.Context, req models.ClassroomRequest) (int32, error) {
	if err := models.ValidateClassroomRequest(req); err != nil {
		return 0, fmt.Errorf("classroom_service.CreateClassroom: %w: %w", errs.ErrBadRequestBody, err)
	}
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("classroom_service.CreateClassroom: %w", err)
	}
	return id, nil
}

func (s *Service) UpdateClassroom(ctx context.Context, req models.Classroom) (int32, error) {
	if err := models.ValidateClassroom(req); err != nil {
		return 0, fmt.Errorf("classroom_service.UpdateClassroom: %w: %w", errs.ErrBadRequestBody, err)
	}
	id, err := s.repo.Update(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("classroom_service.UpdateClassroom: %w", err)
	}
	return id, nil
}

func (s *Service) DeleteClassroom(ctx context.Context, id int32) error {
	if err := models.ValidateClassroomID(id); err != nil {
		return fmt.Errorf("classroom_service.DeleteClassroom: %w: %w", errs.ErrBadRequest, err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("classroom_service.DeleteClassroom: %w", err)
	}
	return nil
}

// schedule
func (s *Service) GetScheduleByClassroomId(ctx context.Context, classroomId int32) ([]*models.Schedule, error) {
	if err := models.ValidateClassroomID(classroomId); err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByClassroomId: %w: %w", errs.ErrBadRequest, err)
	}
	out, err := s.repo.GetScheduleByClassroom(ctx, classroomId)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByClassroomId: %w", err)
	}
	return out, nil
}

func (s *Service) GetScheduleByTeacherId(ctx context.Context, teacherId int32) ([]*models.Schedule, error) {
	if err := models.ValidateTeacherID(teacherId); err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByTeacherId: %w: %w", errs.ErrBadRequest, err)
	}
	out, err := s.repo.GetScheduleByTeacher(ctx, teacherId)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByTeacherId: %w", err)
	}
	return out, nil
}

func (s *Service) CreateSchedule(ctx context.Context, req models.ScheduleRequest) (int32, error) {
	if err := models.ValidateScheduleRequest(req); err != nil {
		return 0, fmt.Errorf("classroom_service.CreateSchedule: %w: %w", errs.ErrBadRequestBody, err)
	}
	id, err := s.repo.CreateSchedule(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("classroom_service.CreateSchedule: %w", err)
	}
	return id, nil
}

func (s *Service) UpdateScheduleById(ctx context.Context, req models.Schedule) error {
	if err := models.ValidateSchedule(req); err != nil {
		return fmt.Errorf("classroom_service.UpdateScheduleById: %w: %w", errs.ErrBadRequestBody, err)
	}
	if err := s.repo.UpdateSchedule(ctx, req); err != nil {
		return fmt.Errorf("classroom_service.UpdateScheduleById: %w", err)
	}
	return nil
}

func (s *Service) DeleteScheduleById(ctx context.Context, id int32) error {
	if err := models.ValidateScheduleID(id); err != nil {
		return fmt.Errorf("classroom_service.DeleteScheduleById: %w: %w", errs.ErrBadRequest, err)
	}
	if err := s.repo.DeleteSchedule(ctx, id); err != nil {
		return fmt.Errorf("classroom_service.DeleteScheduleById: %w", err)
	}
	return nil
}
