package models

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Classroom struct {
	ID                int32  `json:"id"`
	GradeNumber       int32  `json:"grade_number"`
	Letter            string `json:"letter,omitempty"`
	HometownTeacherID *int32  `json:"hometown_teacher_id,omitempty"`
	AcademicYear      string `json:"academic_year"`
}

type ClassroomRequest struct {
	GradeNumber       int32  `json:"grade_number"`
	Letter            string `json:"letter,omitempty"`
	HometownTeacherID *int32  `json:"hometown_teacher_id,omitempty"`
	AcademicYear      string `json:"academic_year"`
}

var (
	ErrInvalidClassroomID     = errors.New("invalid classroom id")
	ErrInvalidGradeNumber     = errors.New("invalid grade number")
	ErrInvalidClassroomLetter = errors.New("invalid classroom letter")
	ErrInvalidTeacherID       = errors.New("invalid teacher id")
	ErrInvalidAcademicYear    = errors.New("invalid academic year")
)

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

func ValidateClassroomLetter(letter string) error {
	letter = strings.TrimSpace(letter)

	if letter == "" {
		return ErrInvalidClassroomLetter
	}

	if utf8.RuneCountInString(letter) != 1 {
		return ErrInvalidClassroomLetter
	}

	return nil
}

func ValidateTeacherID(id *int32) error {
	if id == nil {
		return nil
	}

	if *id <= 0 {
		return ErrInvalidTeacherID
	}

	return nil
}

var academicYearRegex = regexp.MustCompile(`^(20[0-9]{2})\.[0-9]$`)

func ValidateAcademicYear(year string) error {
	year = strings.TrimSpace(year)

	if !academicYearRegex.MatchString(year) {
		return ErrInvalidAcademicYear
	}

	// Получаем часть до точки.
	yearNumber, err := strconv.Atoi((year)[:4])
	if err != nil {
		return ErrInvalidAcademicYear
	}

	currentYear := time.Now().Year()

	if yearNumber < 2000 || yearNumber > currentYear+1 {
		return ErrInvalidAcademicYear
	}

	return nil
}

func ValidateClassroom(classroom Classroom) error {
	if err := ValidateClassroomID(classroom.ID); err != nil {
		return err
	}

	if err := ValidateGradeNumber(classroom.GradeNumber); err != nil {
		return err
	}

	if err := ValidateClassroomLetter(classroom.Letter); err != nil {
		return err
	}

	if err := ValidateTeacherID(classroom.HometownTeacherID); err != nil {
		return err
	}

	if err := ValidateAcademicYear(classroom.AcademicYear); err != nil {
		return err
	}

	return nil
}

func ValidateClassroomRequest(req ClassroomRequest) error {
	if err := ValidateGradeNumber(req.GradeNumber); err != nil {
		return err
	}

	if err := ValidateClassroomLetter(req.Letter); err != nil {
		return err
	}

	if err := ValidateTeacherID(req.HometownTeacherID); err != nil {
		return err
	}

	if err := ValidateAcademicYear(req.AcademicYear); err != nil {
		return err
	}

	return nil
}

func ClassroomName(classroom Classroom) string {
	if strings.TrimSpace(classroom.Letter) == "" {
		return strconv.Itoa(int(classroom.GradeNumber))
	}

	return strconv.Itoa(int(classroom.GradeNumber)) +
		strings.TrimSpace(classroom.Letter)
}
