package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abozorov/school_online/cmd/raiting/internal/models"
	raitingv1 "github.com/abozorov/school_online/grpc_api/generate/raitingpb/raiting/v1"
)

type Service struct {
	repo models.RaitingRepository
}

func New(repo models.RaitingRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func parseDateRange(period string) (string, string) {
	if period == "" {
		// default to current month
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, -1)
		return start.Format("02.01.2006"), end.Format("02.01.2006")
	}
	parts := strings.Split(period, "-")
	if len(parts) != 2 {
		return "01.01.1970", time.Now().Format("02.01.2006")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func (s *Service) GetStudentJournal(ctx context.Context, studentId int, dateRange string) (*raitingv1.GetStudentJournalResponse, error) {
	start, end := parseDateRange(dateRange)
	entries, err := s.repo.GetStudentJournal(ctx, int32(studentId), start, end)
	if err != nil {
		return nil, fmt.Errorf("service.GetStudentJournal: %w", err)
	}

	res := &raitingv1.GetStudentJournalResponse{
		StudentJournal: make(map[string]*raitingv1.StudentGrades),
	}

	for _, e := range entries {
		dateKey := e.Date.Format("02.01.2006")
		if _, ok := res.StudentJournal[dateKey]; !ok {
			res.StudentJournal[dateKey] = &raitingv1.StudentGrades{Grades: make(map[string]string)}
		}
		subjectKey := strconv.Itoa(int(e.SubjectID))
		// decide value: prefer grade, else attendance, else homework
		var val string
		if e.Grade != nil {
			val = strconv.Itoa(int(*e.Grade))
		} else if e.Attendance != nil {
			if *e.Attendance {
				val = "present"
			} else {
				val = "absent"
			}
		} else if e.Homework != nil {
			val = *e.Homework
		} else {
			val = ""
		}
		res.StudentJournal[dateKey].Grades[subjectKey] = val
	}

	return res, nil
}

func (s *Service) GetClassroomJournal(ctx context.Context, classroomId int, dateRange string) (*raitingv1.GetClassroomJournalResponse, error) {
	start, end := parseDateRange(dateRange)
	entries, err := s.repo.GetClassroomJournal(ctx, int32(classroomId), start, end)
	if err != nil {
		return nil, fmt.Errorf("service.GetClassroomJournal: %w", err)
	}

	res := &raitingv1.GetClassroomJournalResponse{
		ClassroomJournal: make(map[string]*raitingv1.Students),
	}

	for _, e := range entries {
		dateKey := e.Date.Format("02.01.2006")
		if _, ok := res.ClassroomJournal[dateKey]; !ok {
			res.ClassroomJournal[dateKey] = &raitingv1.Students{Grades: make(map[string]*raitingv1.StudentGrades)}
		}
		studentKey := strconv.Itoa(int(e.StudentID))
		if _, ok := res.ClassroomJournal[dateKey].Grades[studentKey]; !ok {
			res.ClassroomJournal[dateKey].Grades[studentKey] = &raitingv1.StudentGrades{Grades: make(map[string]string)}
		}
		subjectKey := strconv.Itoa(int(e.SubjectID))
		var val string
		if e.Grade != nil {
			val = strconv.Itoa(int(*e.Grade))
		} else if e.Attendance != nil {
			if *e.Attendance {
				val = "present"
			} else {
				val = "absent"
			}
		} else if e.Homework != nil {
			val = *e.Homework
		} else {
			val = ""
		}
		res.ClassroomJournal[dateKey].Grades[studentKey].Grades[subjectKey] = val
	}

	return res, nil
}

func (s *Service) UpdateGrade(ctx context.Context, req *raitingv1.UpdateJournalRequest) (*raitingv1.UpdateJournalResponse, error) {
	// basic validation
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}
	// parse date
	if req.Date == "" {
		return nil, fmt.Errorf("empty date")
	}
	// build model
	j := models.Journal{
		ClassroomID:  req.ClassId,
		SubjectID:    req.SubjectId,
		TeacherID:    req.TeacherId,
		LessonNumber: req.LessonNumber,
		StudentID:    req.StudentId,
	}
	if req.Attendance != nil {
		j.Attendance = req.Attendance
	}
	if req.Grade != nil {
		j.Grade = req.Grade
	}
	if req.Homework != nil {
		j.Homework = req.Homework
	}

	if err := s.repo.UpsertJournal(ctx, j, req.Date); err != nil {
		return nil, fmt.Errorf("service.UpdateGrade: %w", err)
	}
	return &raitingv1.UpdateJournalResponse{}, nil
}
