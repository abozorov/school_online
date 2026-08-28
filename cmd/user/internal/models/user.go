package models

import (
	"context"
	"errors"
)

type UserService interface {
	Get(ctx context.Context, id int32) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (int32, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int32) error
	CreateSubject(ctx context.Context, name string, description string) (int32, error)
	GetSubjectById(ctx context.Context, id int32) (*Subject, error)
	GetAllSubjects(ctx context.Context) ([]*Subject, error)
	UpdateSubject(ctx context.Context, id int32, name string, description string) error
}

type UserRepository interface {
	Get(ctx context.Context, id int32) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (int32, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int32) error
	CreateSubject(ctx context.Context, name string, description string) (int32, error)
	GetSubjectById(ctx context.Context, id int32) (*Subject, error)
	GetAllSubjects(ctx context.Context) ([]*Subject, error)
	UpdateSubject(ctx context.Context, id int32, name string, description string) error
}

type User struct {
	ID           int32
	Name         string
	Username     string
	Email        string
	VerifyEmail  bool
	PhoneNumber  string
	PasswordHash string
	RefreshToken string
	Role         string
	Birthday     string

	StudentRole *StudentRole
	ParentRole  *ParentRole
	StaffRole   *StaffRole
	TeacherRole *TeacherRole
}

type StudentRole struct {
	ClassroomID int32
}

type ParentRole struct {
	StudentsID []int32
}

type StaffRole struct {
	Position   string
	Experience int32
}

type TeacherRole struct {
	SubjectsID []int32
	Experience int32
}

type Subject struct {
	ID          int32
	Name        string
	Description string
}

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidBirthday    = errors.New("invalid birthday")
	ErrInvalidRole        = errors.New("invalid role")

	ErrEmptyName     = errors.New("empty name")
	ErrEmptyEmail    = errors.New("empty email")
	ErrEmptyPhone    = errors.New("empty phone")
	ErrEmptyUserID   = errors.New("empty user id")
	ErrInvalidAge    = errors.New("invalid age")
	ErrInvalidPhone  = errors.New("invalid phone")
	ErrInvalidUSerId = errors.New("invalid user id")
)
