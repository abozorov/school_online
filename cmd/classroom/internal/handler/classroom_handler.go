package handler

import (
	"context"

	"github.com/abozorov/school_online/cmd/classroom/internal/models"
	classroomv1 "github.com/abozorov/school_online/grpc_api/generate/classroompb/classroom/v1"
	"github.com/abozorov/school_online/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	classroomv1.UnimplementedClassroomServiceServer
	logger  *logger.Logger
	service models.ClassroomService
}

func New(logger *logger.Logger, service models.ClassroomService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// GetClassroom implements ClassroomServiceServer
func (h *Handler) GetClassroom(ctx context.Context, req *classroomv1.GetClassroomRequest) (*classroomv1.ClassroomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	c, err := h.service.GetClassroom(ctx, req.GetId())
	if err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.ClassroomResponse{
		Id:                c.ID,
		GradeNumber:       c.GradeNumber,
		Letter:            c.Letter,
		HometownTeacherId: c.HometownTeacherID,
		AcademicYear:      c.AcademicYear,
	}, nil
}

func (h *Handler) ListClassrooms(ctx context.Context, req *classroomv1.ListClassroomsRequest) (*classroomv1.ListClassroomsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	list, err := h.service.ListClassrooms(ctx)
	if err != nil {
		return nil, responseErr(err)
	}
	resp := &classroomv1.ListClassroomsResponse{Classrooms: make([]*classroomv1.ClassroomResponse, 0, len(list))}
	for _, c := range list {
		resp.Classrooms = append(resp.Classrooms, &classroomv1.ClassroomResponse{
			Id:                c.ID,
			GradeNumber:       c.GradeNumber,
			Letter:            c.Letter,
			HometownTeacherId: c.HometownTeacherID,
			AcademicYear:      c.AcademicYear,
		})
	}
	return resp, nil
}

func (h *Handler) CreateClassroom(ctx context.Context, req *classroomv1.CreateClassroomRequest) (*classroomv1.CreateClassroomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	id, err := h.service.CreateClassroom(ctx, models.ClassroomRequest{
		GradeNumber:       req.GetGradeNumber(),
		Letter:            req.GetLetter(),
		HometownTeacherID: req.GetHometownTeacherId(),
		AcademicYear:      req.GetAcademicYear(),
	})
	if err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.CreateClassroomResponse{Id: id}, nil
}

func (h *Handler) UpdateClassroom(ctx context.Context, req *classroomv1.UpdateClassroomRequest) (*classroomv1.UpdateClassroomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	c := models.Classroom{ID: req.GetId()}
	if req.GradeNumber != nil {
		c.GradeNumber = req.GetGradeNumber()
	}
	if req.Letter != nil {
		c.Letter = req.GetLetter()
	}
	if req.HometownTeacherId != nil {
		c.HometownTeacherID = req.GetHometownTeacherId()
	}
	if req.AcademicYear != nil {
		c.AcademicYear = req.GetAcademicYear()
	}
	id, err := h.service.UpdateClassroom(ctx, c)
	if err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.UpdateClassroomResponse{Id: id}, nil
}

func (h *Handler) DeleteClassroom(ctx context.Context, req *classroomv1.DeleteClassroomRequest) (*classroomv1.DeleteClassroomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	if err := h.service.DeleteClassroom(ctx, req.GetId()); err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.DeleteClassroomResponse{}, nil
}

func (h *Handler) GetScheduleByClassroom(ctx context.Context, req *classroomv1.GetScheduleByClassroomRequest) (*classroomv1.ScheduleList, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	list, err := h.service.GetScheduleByClassroomId(ctx, req.GetClassId())
	if err != nil {
		return nil, responseErr(err)
	}
	out := &classroomv1.ScheduleList{Schedules: make([]*classroomv1.ScheduleResponse, 0, len(list))}
	for _, s := range list {
		out.Schedules = append(out.Schedules, &classroomv1.ScheduleResponse{
			Id:           s.ID,
			ClassId:      s.ClassroomID,
			SubjectId:    s.SubjectID,
			TeacherId:    s.TeacherID,
			DayOfWeek:    s.DayOfWeek,
			LessonNumber: s.LessonNumber,
			Room:         s.Room,
			AcademicYear: s.AcademicYear,
		})
	}
	return out, nil
}

func (h *Handler) GetScheduleByTeacher(ctx context.Context, req *classroomv1.GetScheduleByTeacherRequest) (*classroomv1.ScheduleList, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	list, err := h.service.GetScheduleByTeacherId(ctx, req.GetTeacherId())
	if err != nil {
		return nil, responseErr(err)
	}
	out := &classroomv1.ScheduleList{Schedules: make([]*classroomv1.ScheduleResponse, 0, len(list))}
	for _, s := range list {
		out.Schedules = append(out.Schedules, &classroomv1.ScheduleResponse{
			Id:           s.ID,
			ClassId:      s.ClassroomID,
			SubjectId:    s.SubjectID,
			TeacherId:    s.TeacherID,
			DayOfWeek:    s.DayOfWeek,
			LessonNumber: s.LessonNumber,
			Room:         s.Room,
			AcademicYear: s.AcademicYear,
		})
	}
	return out, nil
}

func (h *Handler) CreateSchedule(ctx context.Context, req *classroomv1.CreateScheduleRequest) (*classroomv1.CreateScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	id, err := h.service.CreateSchedule(ctx, models.ScheduleRequest{
		ClassroomID:  req.GetClassId(),
		SubjectID:    req.GetSubjectId(),
		TeacherID:    req.GetTeacherId(),
		DayOfWeek:    req.GetDayOfWeek(),
		LessonNumber: req.GetLessonNumber(),
		Room:         req.GetRoom(),
		AcademicYear: req.GetAcademicYear(),
	})
	if err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.CreateScheduleResponse{ClassId: id}, nil
}

func (h *Handler) UpdateSchedule(ctx context.Context, req *classroomv1.UpdateScheduleRequest) (*classroomv1.UpdateScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	s := models.Schedule{ID: req.GetId()}
	if req.ClassId != nil {
		s.ClassroomID = req.GetClassId()
	}
	if req.SubjectId != nil {
		s.SubjectID = req.GetSubjectId()
	}
	if req.TeacherId != nil {
		s.TeacherID = req.GetTeacherId()
	}
	if req.DayOfWeek != nil {
		s.DayOfWeek = req.GetDayOfWeek()
	}
	if req.LessonNumber != nil {
		s.LessonNumber = req.GetLessonNumber()
	}
	if req.Room != nil {
		s.Room = req.GetRoom()
	}
	if req.AcademicYear != nil {
		s.AcademicYear = req.GetAcademicYear()
	}
	if err := h.service.UpdateScheduleById(ctx, s); err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.UpdateScheduleResponse{ClassId: s.ClassroomID}, nil
}

func (h *Handler) DeleteSchedule(ctx context.Context, req *classroomv1.DeleteScheduleRequest) (*classroomv1.DeleteScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Classroom microservice: request is required")
	}
	if err := h.service.DeleteScheduleById(ctx, req.GetId()); err != nil {
		return nil, responseErr(err)
	}
	return &classroomv1.DeleteScheduleResponse{}, nil
}
