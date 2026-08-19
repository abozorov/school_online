package raitingservice

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/api_gateway/internal/models"
	"github.com/abozorov/school_online/cmd/api_gateway/internal/services"
	raitingv1 "github.com/abozorov/school_online/grpc_api/generate/raitingpb/raiting/v1"
)

type RaitingService struct {
	serviceManager services.IServiceManager
}

func NewRaitingService(serviceManager services.IServiceManager) *RaitingService {

	return &RaitingService{
		serviceManager: serviceManager,
	}
}

func (r *RaitingService) GetJournalByStudentId(ctx context.Context, studentId int, dateRange string) (map[string]map[string]string, error) {
	// get journal by student id
	journal, err := r.serviceManager.RaitingService().GetStudentJournal(ctx, &raitingv1.GetJournalRequest{
		Id:     int32(studentId),
		Period: dateRange,
	})
	if err != nil {
		return nil, fmt.Errorf("raiting_service.GetJournalByStudentId: %w", err)
	}

	out := make(map[string]map[string]string)
	for date, studenJournal := range journal.StudentJournal {
		out[date] = make(map[string]string)
		for subject, studentGrade := range studenJournal.Grades {
			out[date][subject] = studentGrade
		}
	}

	return out, nil
}

func (r *RaitingService) GetJournalByClassroomId(ctx context.Context, classroomId int, dateRange string) (map[string]map[string]map[string]string, error) {
	// get journal by classroom id
	journal, err := r.serviceManager.RaitingService().GetClassroomJournal(ctx, &raitingv1.GetJournalRequest{
		Id:     int32(classroomId),
		Period: dateRange,
	})
	if err != nil {
		return nil, fmt.Errorf("raiting_service.GetJournalByClassroomId: %w", err)
	}

	out := make(map[string]map[string]map[string]string)
	for date, students := range journal.ClassroomJournal {
		out[date] = make(map[string]map[string]string)
		for student, subjects := range students.Grades {
			out[date][student] = make(map[string]string)
			for subject, grade := range subjects.Grades {
				out[date][student][subject] = grade
			}
		}
	}

	return out, nil
}

func (r *RaitingService) UpdateJournal(ctx context.Context, req models.Journal) error {
	err := models.ValidateJournal(req)
	if err != nil {
		return fmt.Errorf("raiting_service.UpdateJournal: %w", err)
	}

	_, err = r.serviceManager.RaitingService().UpdateGrade(ctx, &raitingv1.UpdateJournalRequest{
		ClassId:    req.ClassroomID,
		SubjectId:  req.SubjectID,
		Date:       req.Date.Format("02.01.2026"),
		StudentId:  req.StudentID,
		Attendance: &req.Attendance,
		Grade:      &req.Grade,
		Homework:   &req.Homework,
	})
	if err != nil {
		return fmt.Errorf("raiting_service.UpdateJournal: %w", err)
	}

	return nil
}
