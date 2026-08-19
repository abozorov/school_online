package models

import (
	"errors"
	"time"
)

type Journal struct {
	ClassroomID  int32     `json:"classroom_id"`
	SubjectID    int32     `json:"subject_id"`
	TeacherID    int32     `json:"teacher_id"`
	Date         time.Time `json:"date"`
	LessonNumber int32     `json:"lesson_number"`
	Attendance   *bool     `json:"attendance,omitempty"`
	StudentID    int32     `json:"student_id"`
	Grade        *int32    `json:"grade,omitempty"`
	Homework     *string   `json:"homework,omitempty"`
}

var (
	ErrInvalidJournalID = errors.New("invalid journal id")
	ErrInvalidDate      = errors.New("invalid date")
	ErrInvalidStudentID = errors.New("invalid student id")
	ErrInvalidGrade     = errors.New("invalid grade")
	ErrInvalidHomework  = errors.New("invalid homework")
)

func ValidateJournalID(id int32) error {
	if id <= 0 {
		return ErrInvalidJournalID
	}

	return nil
}

func ValidateDate(date time.Time) error {
	if date.IsZero() {
		return ErrInvalidDate
	}

	// Дата не может быть в будущем.
	if date.After(time.Now()) {
		return ErrInvalidDate
	}

	return nil
}

func ValidateStudentID(id int32) error {
	if id <= 0 {
		return ErrInvalidStudentID
	}

	return nil
}

func ValidateGrade(grade *int32) error {
	if grade == nil {
		return nil
	}

	if *grade < 1 || *grade > 5 {
		return ErrInvalidGrade
	}

	return nil
}

func ValidateHomework(homework *string) error {
	if homework == nil {
		return nil
	}

	if len(*homework) > 5000 {
		return ErrInvalidHomework
	}

	return nil
}

func ValidateJournal(journal Journal) error {
	if err := ValidateClassroomID(journal.ClassroomID); err != nil {
		return err
	}

	if err := ValidateSubjectID(journal.SubjectID); err != nil {
		return err
	}

	if err := ValidateTeacherID(&journal.TeacherID); err != nil {
		return err
	}

	if err := ValidateDate(journal.Date); err != nil {
		return err
	}

	if err := ValidateLessonNumber(journal.LessonNumber); err != nil {
		return err
	}

	if err := ValidateStudentID(journal.StudentID); err != nil {
		return err
	}

	if err := ValidateGrade(journal.Grade); err != nil {
		return err
	}

	if err := ValidateHomework(journal.Homework); err != nil {
		return err
	}

	return nil
}
