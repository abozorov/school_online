package models

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Domain models
type Classroom struct {
	ID                int32  `json:"id"`
	GradeNumber       int32  `json:"grade_number"`
	Letter            string `json:"letter,omitempty"`
	HometownTeacherID int32  `json:"hometown_teacher_id,omitempty"`
	AcademicYear      string `json:"academic_year"`
}

type ClassroomRequest struct {
	GradeNumber       int32  `json:"grade_number"`
	Letter            string `json:"letter,omitempty"`
	HometownTeacherID int32  `json:"hometown_teacher_id,omitempty"`
	AcademicYear      string `json:"academic_year"`
}

type Schedule struct {
	ID           int32  `json:"id"`
	ClassroomID  int32  `json:"classroom_id"`
	SubjectID    int32  `json:"subject_id"`
	TeacherID    int32  `json:"teacher_id"`
	DayOfWeek    int32  `json:"day_of_week"`
	LessonNumber int32  `json:"lesson_number"`
	Room         int32  `json:"room,omitempty"`
	AcademicYear string `json:"academic_year"`
}

type ScheduleRequest struct {
	ClassroomID  int32  `json:"classroom_id"`
	SubjectID    int32  `json:"subject_id"`
	TeacherID    int32  `json:"teacher_id"`
	DayOfWeek    int32  `json:"day_of_week"`
	LessonNumber int32  `json:"lesson_number"`
	Room         int32  `json:"room,omitempty"`
	AcademicYear string `json:"academic_year"`
}

// Errors
var (
	ErrInvalidClassroomID     = errors.New("invalid classroom id")
	ErrInvalidGradeNumber     = errors.New("invalid grade number")
	ErrInvalidClassroomLetter = errors.New("invalid classroom letter")
	ErrInvalidTeacherID       = errors.New("invalid teacher id")
	ErrInvalidAcademicYear    = errors.New("invalid academic year")

	ErrInvalidScheduleID   = errors.New("invalid schedule id")
	ErrInvalidSubjectID    = errors.New("invalid subject id")
	ErrInvalidDayOfWeek    = errors.New("invalid day of week")
	ErrInvalidLessonNumber = errors.New("invalid lesson number")
	ErrInvalidRoom         = errors.New("invalid room")
)

// Validators
func ValidateClassroomID(id int32) error {
	if id <= 0 {
		return ErrInvalidClassroomID
	}
	return nil
}

func ValidateGradeNumber(grade int32) error {
	if grade < 1 || grade > 11 {
		return ErrInvalidGradeNumber
	}
	return nil
}

func ValidateClassroomLetter(letter *string) error {
	*letter = strings.TrimSpace(*letter)
	if *letter == "" {
		return ErrInvalidClassroomLetter
	}
	if len([]rune(*letter)) != 1 {
		return ErrInvalidClassroomLetter
	}
	return nil
}

func ValidateTeacherID(id int32) error {
	if id <= 0 {
		return ErrInvalidTeacherID
	}
	return nil
}

var academicYearRegex = regexp.MustCompile(`^(20[0-9]{2})\.[0-9]$`)

func ValidateAcademicYear(year *string) error {
	*year = strings.TrimSpace(*year)
	if !academicYearRegex.MatchString(*year) {
		return ErrInvalidAcademicYear
	}
	yearNumber, err := strconv.Atoi((*year)[:4])
	if err != nil {
		return ErrInvalidAcademicYear
	}
	currentYear := time.Now().Year()
	if yearNumber < 2000 || yearNumber > currentYear+1 {
		return ErrInvalidAcademicYear
	}
	return nil
}

func ValidateClassroom(c Classroom) error {
	if err := ValidateClassroomID(c.ID); err != nil {
		return err
	}
	if err := ValidateGradeNumber(c.GradeNumber); err != nil {
		return err
	}
	if err := ValidateClassroomLetter(&c.Letter); err != nil {
		return err
	}
	if c.HometownTeacherID != 0 {
		if err := ValidateTeacherID(c.HometownTeacherID); err != nil {
			return err
		}
	}
	if err := ValidateAcademicYear(&c.AcademicYear); err != nil {
		return err
	}
	return nil
}

func ValidateClassroomRequest(req ClassroomRequest) error {
	if err := ValidateGradeNumber(req.GradeNumber); err != nil {
		return err
	}
	if err := ValidateClassroomLetter(&req.Letter); err != nil {
		return err
	}
	if req.HometownTeacherID != 0 {
		if err := ValidateTeacherID(req.HometownTeacherID); err != nil {
			return err
		}
	}
	if err := ValidateAcademicYear(&req.AcademicYear); err != nil {
		return err
	}
	return nil
}

// Schedule validators
func ValidateScheduleID(id int32) error {
	if id <= 0 {
		return ErrInvalidScheduleID
	}
	return nil
}

func ValidateSubjectID(id int32) error {
	if id <= 0 {
		return ErrInvalidSubjectID
	}
	return nil
}

func ValidateDayOfWeek(day int32) error {
	if day < 1 || day > 7 {
		return ErrInvalidDayOfWeek
	}
	return nil
}

func ValidateLessonNumber(number int32) error {
	if number < 1 || number > 8 {
		return ErrInvalidLessonNumber
	}
	return nil
}

func ValidateRoom(room int32) error {
	if room < 0 {
		return ErrInvalidRoom
	}
	return nil
}

func ValidateSchedule(s Schedule) error {
	if err := ValidateScheduleID(s.ID); err != nil {
		return err
	}
	if err := ValidateClassroomID(s.ClassroomID); err != nil {
		return err
	}
	if err := ValidateSubjectID(s.SubjectID); err != nil {
		return err
	}
	if err := ValidateTeacherID(s.TeacherID); err != nil {
		return err
	}
	if err := ValidateDayOfWeek(s.DayOfWeek); err != nil {
		return err
	}
	if err := ValidateLessonNumber(s.LessonNumber); err != nil {
		return err
	}
	if err := ValidateRoom(s.Room); err != nil {
		return err
	}
	if err := ValidateAcademicYear(&s.AcademicYear); err != nil {
		return err
	}
	return nil
}

func ValidateScheduleRequest(req ScheduleRequest) error {
	if err := ValidateClassroomID(req.ClassroomID); err != nil {
		return err
	}
	if err := ValidateSubjectID(req.SubjectID); err != nil {
		return err
	}
	if err := ValidateTeacherID(req.TeacherID); err != nil {
		return err
	}
	if err := ValidateDayOfWeek(req.DayOfWeek); err != nil {
		return err
	}
	if err := ValidateLessonNumber(req.LessonNumber); err != nil {
		return err
	}
	if err := ValidateRoom(req.Room); err != nil {
		return err
	}
	if err := ValidateAcademicYear(&req.AcademicYear); err != nil {
		return err
	}
	return nil
}

// Interfaces
type ClassroomService interface {
	GetClassroom(ctx context.Context, id int32) (*Classroom, error)
	ListClassrooms(ctx context.Context) ([]*Classroom, error)
	CreateClassroom(ctx context.Context, req ClassroomRequest) (int32, error)
	UpdateClassroom(ctx context.Context, req Classroom) (int32, error)
	DeleteClassroom(ctx context.Context, id int32) error

	GetScheduleByClassroomId(ctx context.Context, classroomId int32) ([]*Schedule, error)
	GetScheduleByTeacherId(ctx context.Context, teacherId int32) ([]*Schedule, error)
	CreateSchedule(ctx context.Context, req ScheduleRequest) (int32, error)
	UpdateScheduleById(ctx context.Context, req Schedule) error
	DeleteScheduleById(ctx context.Context, id int32) error
}

type ClassroomRepository interface {
	Get(ctx context.Context, id int32) (*Classroom, error)
	List(ctx context.Context) ([]*Classroom, error)
	Create(ctx context.Context, req ClassroomRequest) (int32, error)
	Update(ctx context.Context, req Classroom) (int32, error)
	Delete(ctx context.Context, id int32) error

	GetScheduleByClassroom(ctx context.Context, classroomId int32) ([]*Schedule, error)
	GetScheduleByTeacher(ctx context.Context, teacherId int32) ([]*Schedule, error)
	CreateSchedule(ctx context.Context, req ScheduleRequest) (int32, error)
	UpdateSchedule(ctx context.Context, req Schedule) error
	DeleteSchedule(ctx context.Context, id int32) error
}
