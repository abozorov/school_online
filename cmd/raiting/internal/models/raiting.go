package models

import (
	"context"
	"errors"
	"time"

	raitingv1 "github.com/abozorov/school_online/grpc_api/generate/raitingpb/raiting/v1"
)

// Journal represents a journal entry row
type Journal struct {
	ID           int32      `json:"id"`
	ClassroomID  int32      `json:"classroom_id"`
	SubjectID    int32      `json:"subject_id"`
	TeacherID    int32      `json:"teacher_id"`
	Date         time.Time  `json:"date"`
	LessonNumber int32      `json:"lesson_number"`
	Room         int32      `json:"room,omitempty"`
	Attendance   *bool      `json:"attendance,omitempty"`
	StudentID    int32      `json:"student_id"`
	Grade        *int32     `json:"grade,omitempty"`
	Homework     *string    `json:"homework,omitempty"`
}

var (
	ErrInvalidJournalID = errors.New("invalid journal id")
)

func ValidateJournalID(id int32) error {
	if id <= 0 {
		return ErrInvalidJournalID
	}
	return nil
}

// Interfaces
type RaitingRepository interface {
	GetStudentJournal(ctx context.Context, studentId int32, startDate, endDate string) ([]Journal, error)
	GetClassroomJournal(ctx context.Context, classroomId int32, startDate, endDate string) ([]Journal, error)
	UpsertJournal(ctx context.Context, j Journal, dateStr string) error
}

type RaitingService interface {
	GetStudentJournal(ctx context.Context, studentId int, dateRange string) (*raitingv1.GetStudentJournalResponse, error)
	GetClassroomJournal(ctx context.Context, classroomId int, dateRange string) (*raitingv1.GetClassroomJournalResponse, error)
	UpdateGrade(ctx context.Context, req *raitingv1.UpdateJournalRequest) (*raitingv1.UpdateJournalResponse, error)
}
