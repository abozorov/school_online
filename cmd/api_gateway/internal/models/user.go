package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/abozorov/school_online/pkg/errs"
)

type User struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	VerifyEmail  bool   `json:"verify_email"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	PasswordHash string `json:"-"`
	RefreshToken string `json:"-"`
	Role         string `json:"role,omitempty"`
	Birthday     string `json:"birthday,omitempty"`

	StudentRole *StudentRole `json:"student_role,omitempty"`
	ParentRole  *ParentRole  `json:"parent_role,omitempty"`
	StaffRole   *StaffRole   `json:"staff_role,omitempty"`
	TeacherRole *TeacherRole `json:"teacher_role,omitempty"`
}

type RegisterUserRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty"`
	Password    string `json:"password"`
	Birthday    string `json:"birthday,omitempty"`

	StudentRole *StudentRole `json:"student_role,omitempty"`
	ParentRole  *ParentRole  `json:"parent_role,omitempty"`
	StaffRole   *StaffRole   `json:"staff_role,omitempty"`
	TeacherRole *TeacherRole `json:"teacher_role,omitempty"`
}

type UpdateUserRequest struct {
	ID          int32  `json:"id"`
	Name        string `json:"name,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty"`
	Birthday    string `json:"birthday,omitempty"`

	StudentRole *StudentRole `json:"student_role,omitempty"`
	ParentRole  *ParentRole  `json:"parent_role,omitempty"`
	StaffRole   *StaffRole   `json:"staff_role,omitempty"`
	TeacherRole *TeacherRole `json:"teacher_role,omitempty"`
}

type StudentRole struct {
	ClassroomID int32 `json:"classroom_id,omitempty"`
}

type ParentRole struct {
	StudentsID []int32 `json:"students_id"`
}

type StaffRole struct {
	Position   string `json:"position,omitempty"`
	Experience int32  `json:"experience"`
}

type TeacherRole struct {
	SubjectsID []int32 `json:"subjects_id"`
	Experience int32   `json:"experience"`
}

var (
	errInvalidID          = errors.New("invalid id")
	errInvalidName        = errors.New("invalid name")
	errInvalidUsername    = errors.New("invalid username")
	errInvalidEmail       = errors.New("invalid email")
	errInvalidPhoneNumber = errors.New("invalid phone number")
	errInvalidPassword    = errors.New("invalid password")
	errInvalidBirthday    = errors.New("invalid birthday")
	errInvalidRole        = errors.New("invalid role")
)

func ValidateID(id int32) error {
	if id <= 0 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidID)
	}

	return nil
}

func ValidateName(name *string) error {
	*name = strings.TrimSpace(*name)

	if *name == "" {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidName)
	}

	length := utf8.RuneCountInString(*name)

	if length < 2 || length > 100 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidName)
	}

	return nil
}

func ValidateUsername(username *string) error {
	*username = strings.TrimSpace(*username)

	if len(*username) < 3 || len(*username) > 100 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidUsername)
	}

	for _, r := range *username {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_') {
			return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidUsername)
		}
	}

	return nil
}

func ValidateEmail(email *string) error {
	*email = strings.TrimSpace(*email)

	if *email == "" {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidEmail)
	}

	if len(*email) > 100 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidEmail)
	}

	_, err := mail.ParseAddress(*email)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidEmail)
	}

	return nil
}

func ValidatePhoneNumber(phone *string) error {
	if phone == nil {
		return nil
	}

	if len(*phone) < 10 || len(*phone) > 12 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidPhoneNumber)
	}

	for _, r := range *phone {
		if r < '0' || r > '9' {
			return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidPhoneNumber)
		}
	}

	return nil
}

func ValidatePassword(password string) error {
	// bcrypt has a 72-byte limit.
	if len(password) < 8 || len(password) > 72 {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidPassword)
	}

	var hasLetter bool
	var hasDigit bool

	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}

		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidPassword)
	}

	return nil
}

func ValidateBirthday(birthday *string) error {

	birthDate, err := time.Parse("02-01-2006", *birthday)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidBirthday)
	}

	now := time.Now()

	// Birthday cannot be in the future.
	if birthDate.After(now) {
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidBirthday)
	}

	return nil
}

func ValidateRole(role *string) error {
	if role == nil {
		return nil
	}

	switch *role {
	case "user",
		"staff",
		"teacher",
		"student",
		"parent":
		return nil

	default:
		return fmt.Errorf("%w: %w", errs.ErrBadRequestBody, errInvalidRole)
	}
}

func ValidateRegisterRequest(req RegisterUserRequest) error {
	if err := ValidateName(&req.Name); err != nil {
		return err
	}

	if err := ValidateUsername(&req.Username); err != nil {
		return err
	}

	if err := ValidateEmail(&req.Email); err != nil {
		return err
	}

	if err := ValidatePhoneNumber(&req.PhoneNumber); err != nil {
		return err
	}

	if err := ValidatePassword(req.Password); err != nil {
		return err
	}

	if err := ValidateBirthday(&req.Birthday); err != nil {
		return err
	}

	return nil
}

func ValidateUpdateUserRequest(req UpdateUserRequest) error {

	if err := ValidateName(&req.Name); err != nil {
		return err
	}

	if err := ValidateUsername(&req.Username); err != nil {
		return err
	}

	if err := ValidateEmail(&req.Email); err != nil {
		return err
	}

	if err := ValidatePhoneNumber(&req.PhoneNumber); err != nil {
		return err
	}

	if err := ValidateBirthday(&req.Birthday); err != nil {
		return err
	}

	return nil
}
