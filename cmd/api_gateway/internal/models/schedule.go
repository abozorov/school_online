package models

import "errors"

type Schedule struct {
	ID           int32  `json:"id"`
	ClassroomID  int32  `json:"classroom_id"`
	SubjectID    int32  `json:"subject_id"`
	TeacherID    int32  `json:"teacher_id"`
	DayOfWeek    int32  `json:"day_of_week"`
	LessonNumber int32  `json:"lesson_number"`
	Room         *int32  `json:"room,omitempty"`
	AcademicYear string `json:"academic_year"`
}

type ScheduleRequest struct {
	ClassroomID  int32  `json:"classroom_id"`
	SubjectID    int32  `json:"subject_id"`
	TeacherID    int32  `json:"teacher_id"`
	DayOfWeek    int32  `json:"day_of_week"`
	LessonNumber int32  `json:"lesson_number"`
	Room         *int32  `json:"room,omitempty"`
	AcademicYear string `json:"academic_year"`
}

var (
	ErrInvalidScheduleID   = errors.New("invalid schedule id")
	ErrInvalidSubjectID    = errors.New("invalid subject id")
	ErrInvalidDayOfWeek    = errors.New("invalid day of week")
	ErrInvalidLessonNumber = errors.New("invalid lesson number")
	ErrInvalidRoom         = errors.New("invalid room")
)

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

func ValidateRoom(room *int32) error {
	if room == nil {
		return nil
	}
	if *room <= 0 {
		return ErrInvalidRoom
	}

	return nil
}

func ValidateSchedule(schedule Schedule) error {
	if err := ValidateScheduleID(schedule.ID); err != nil {
		return err
	}

	if err := ValidateClassroomID(schedule.ClassroomID); err != nil {
		return err
	}

	if err := ValidateSubjectID(schedule.SubjectID); err != nil {
		return err
	}

	if err := ValidateTeacherID(&schedule.TeacherID); err != nil {
		return err
	}

	if err := ValidateDayOfWeek(schedule.DayOfWeek); err != nil {
		return err
	}

	if err := ValidateLessonNumber(schedule.LessonNumber); err != nil {
		return err
	}

	if err := ValidateRoom(schedule.Room); err != nil {
		return err
	}

	if err := ValidateAcademicYear(schedule.AcademicYear); err != nil {
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

	if err := ValidateTeacherID(&req.TeacherID); err != nil {
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

	if err := ValidateAcademicYear(req.AcademicYear); err != nil {
		return err
	}

	return nil
}
