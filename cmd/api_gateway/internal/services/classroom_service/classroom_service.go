package classroomservice

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	"github.com/abozorov/school_online/pkg/errs"

	classroomv1 "github.com/abozorov/school_online/grpc_api/generate/classroompb/classroom/v1"
)

type ClassroomService struct {
	serviceManager services.IServiceManager
}

func NewClassroomService(serviceManager services.IServiceManager) *ClassroomService {

	return &ClassroomService{
		serviceManager: serviceManager,
	}
}

func (c *ClassroomService) GetClassroomByID(ctx context.Context, id int32) (*models.Classroom, error) {
	// validate id
	err := models.ValidateID(id)
	if err != nil {
		return &models.Classroom{}, fmt.Errorf("classroom_service.GetClassroomByID: %w: %w", errs.ErrBadRequestBody, services.GRPCToErrs(err))
	}

	// get classroom by id
	classroom, err := c.serviceManager.ClassroomService().GetClassroom(ctx, &classroomv1.GetClassroomRequest{
		Id: id,
	})
	if err != nil {
		return &models.Classroom{}, fmt.Errorf("classroom_service.GetClassroomByID: %w", services.GRPCToErrs(err))
	}

	return &models.Classroom{
		ID:                classroom.GetId(),
		GradeNumber:       classroom.GetGradeNumber(),
		Letter:            classroom.GetLetter(),
		HometownTeacherID: &classroom.HometownTeacherId,
		AcademicYear:      classroom.GetAcademicYear(),
	}, nil
}

func (c *ClassroomService) List(ctx context.Context) ([]*models.Classroom, error) {
	// get list of classrooms
	classrooms, err := c.serviceManager.ClassroomService().ListClassrooms(ctx, &classroomv1.ListClassroomsRequest{})
	if err != nil {
		return nil, fmt.Errorf("classroom_service.List: %w", services.GRPCToErrs(err))
	}

	// map to models.Classroom
	var result []*models.Classroom
	for _, classroom := range classrooms.GetClassrooms() {
		result = append(result, &models.Classroom{
			ID:                classroom.GetId(),
			GradeNumber:       classroom.GetGradeNumber(),
			Letter:            classroom.GetLetter(),
			HometownTeacherID: &classroom.HometownTeacherId,
			AcademicYear:      classroom.GetAcademicYear(),
		})
	}

	return result, nil
}

func (c *ClassroomService) Create(ctx context.Context, request models.ClassroomRequest) (int32, error) {
	// validate request
	err := models.ValidateClassroomRequest(request)
	if err != nil {
		return 0, fmt.Errorf("classroom_service.Create: %w: %w", errs.ErrBadRequestBody, services.GRPCToErrs(err))
	}

	// create classroom
	out, err := c.serviceManager.ClassroomService().CreateClassroom(ctx, &classroomv1.CreateClassroomRequest{
		GradeNumber:       request.GradeNumber,
		Letter:            request.Letter,
		HometownTeacherId: *request.HometownTeacherID,
		AcademicYear:      request.AcademicYear,
	})
	if err != nil {
		return 0, fmt.Errorf("classroom_service.Create: %w", services.GRPCToErrs(err))
	}

	return out.GetId(), nil
}

func (c *ClassroomService) UpdateByID(ctx context.Context, request models.Classroom) error {
	// validate request
	err := models.ValidateClassroom(request)
	if err != nil {
		return fmt.Errorf("classroom_service.UpdateByID: %w: %w", errs.ErrBadRequestBody, services.GRPCToErrs(err))
	}

	// update classroom
	_, err = c.serviceManager.ClassroomService().UpdateClassroom(ctx, &classroomv1.UpdateClassroomRequest{
		Id:                request.ID,
		GradeNumber:       &request.GradeNumber,
		Letter:            &request.Letter,
		HometownTeacherId: request.HometownTeacherID,
		AcademicYear:      &request.AcademicYear,
	})
	if err != nil {
		return fmt.Errorf("classroom_service.UpdateByID: %w", services.GRPCToErrs(err))
	}

	return nil
}

func (c *ClassroomService) DeleteByID(ctx context.Context, id int32) error {
	// validate id
	err := models.ValidateID(id)
	if err != nil {
		return fmt.Errorf("classroom_service.DeleteByID: %w: %w", errs.ErrBadRequest, services.GRPCToErrs(err))
	}

	// delete classroom
	_, err = c.serviceManager.ClassroomService().DeleteClassroom(ctx, &classroomv1.DeleteClassroomRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("classroom_service.DeleteByID: %w", services.GRPCToErrs(err))
	}

	return nil
}