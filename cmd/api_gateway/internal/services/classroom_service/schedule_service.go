package classroomservice

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	classroomv1 "github.com/abozorov/school_online/grpc_api/generate/classroompb/classroom/v1"
	"github.com/abozorov/school_online/pkg/errs"
)

func (c *ClassroomService) GetScheduleByClassroomId(ctx context.Context, classroomId int32) ([]*models.Schedule, error) {
	// validate id
	err := models.ValidateID(classroomId)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByClassroomId: %w: %w", errs.ErrBadRequest, err)
	}

	// get schedule by classroom id
	out, err := c.serviceManager.ClassroomService().GetScheduleByClassroom(ctx, &classroomv1.GetScheduleByClassroomRequest{
		ClassId: classroomId,
	})
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByClassroomId: %w", err)
	}

	schedules := make([]*models.Schedule, len(out.GetSchedules()))
	for i, s := range out.GetSchedules() {
		schedules[i] = &models.Schedule{
			ID:           s.GetId(),
			ClassroomID:  s.GetClassId(),
			SubjectID:    s.GetSubjectId(),
			TeacherID:    s.GetTeacherId(),
			DayOfWeek:    s.GetDayOfWeek(),
			LessonNumber: s.GetLessonNumber(),
			Room:         s.GetRoom(),
			AcademicYear: s.GetAcademicYear(),
		}
	}

	return schedules, nil
}

func (c *ClassroomService) GetScheduleByTeacherId(ctx context.Context, teacherId int32) ([]*models.Schedule, error) {
	// validate id
	err := models.ValidateID(teacherId)
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByTeacherId: %w: %w", errs.ErrBadRequest, err)
	}

	// get schedule by teacher id
	out, err := c.serviceManager.ClassroomService().GetScheduleByTeacher(ctx, &classroomv1.GetScheduleByTeacherRequest{
		TeacherId: teacherId,
	})
	if err != nil {
		return nil, fmt.Errorf("classroom_service.GetScheduleByTeacherId: %w", err)
	}

	schedules := make([]*models.Schedule, len(out.GetSchedules()))
	for i, s := range out.GetSchedules() {
		schedules[i] = &models.Schedule{
			ID:           s.GetId(),
			ClassroomID:  s.GetClassId(),
			SubjectID:    s.GetSubjectId(),
			TeacherID:    s.GetTeacherId(),
			DayOfWeek:    s.GetDayOfWeek(),
			LessonNumber: s.GetLessonNumber(),
			Room:         s.GetRoom(),
			AcademicYear: s.GetAcademicYear(),
		}
	}

	return schedules, nil
}

func (c *ClassroomService) CreateSchedule(ctx context.Context, request models.ScheduleRequest) (int32, error) {
	// validate request
	err := models.ValidateScheduleRequest(request)
	if err != nil {
		return 0, fmt.Errorf("classroom_service.CreateSchedule: %w: %w", errs.ErrBadRequestBody, err)
	}

	// create schedule
	out, err := c.serviceManager.ClassroomService().CreateSchedule(ctx, &classroomv1.CreateScheduleRequest{
		ClassId:      request.ClassroomID,
		SubjectId:    request.SubjectID,
		TeacherId:    request.TeacherID,
		DayOfWeek:    request.DayOfWeek,
		LessonNumber: request.LessonNumber,
		Room:         request.Room,
		AcademicYear: request.AcademicYear,
	})
	if err != nil {
		return 0, fmt.Errorf("classroom_service.CreateSchedule: %w", err)
	}

	return out.GetClassId(), nil
}

func (c *ClassroomService) UpdateScheduleById(ctx context.Context, request models.Schedule) error {
	// validate request
	err := models.ValidateSchedule(request)
	if err != nil {
		return fmt.Errorf("classroom_service.UpdateScheduleById: %w: %w", errs.ErrBadRequestBody, err)
	}

	// update schedule
	_, err = c.serviceManager.ClassroomService().UpdateSchedule(ctx, &classroomv1.UpdateScheduleRequest{
		Id:           request.ID,
		ClassId:      &request.ClassroomID,
		SubjectId:    &request.SubjectID,
		TeacherId:    &request.TeacherID,
		DayOfWeek:    &request.DayOfWeek,
		LessonNumber: &request.LessonNumber,
		Room:         &request.Room,
		AcademicYear: &request.AcademicYear,
	})
	if err != nil {
		return fmt.Errorf("classroom_service.UpdateScheduleById: %w", err)
	}

	return nil
}

func (c *ClassroomService) DeleteScheduleById(ctx context.Context, id int32) error {
	// validate id
	err := models.ValidateID(id)
	if err != nil {
		return fmt.Errorf("classroom_service.DeleteScheduleById: %w: %w", errs.ErrBadRequest, err)
	}

	// delete schedule
	_, err = c.serviceManager.ClassroomService().DeleteSchedule(ctx, &classroomv1.DeleteScheduleRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("classroom_service.DeleteScheduleById: %w", err)
	}

	return nil
}