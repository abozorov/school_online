package models

import (
	"strings"

	"github.com/abozorov/school_online/pkg/errs"
)

type Subject struct {
	ID          int32  `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func ValidateSubjectName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errs.ErrBadRequestBody
	}
	return nil
}

func ValidateSubjectDescription(description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return errs.ErrBadRequestBody
	}
	return nil
}

func ValidateSubject(s *Subject) error {
	if err := ValidateSubjectName(s.Name); err != nil {
		return err
	}

	if err := ValidateSubjectDescription(s.Description); err != nil {
		return err
	}

	return nil
}

func ValidateCreateSubjectRequest(req CreateSubjectRequest) error {
	if err := ValidateSubjectName(req.Name); err != nil {
		return err
	}

	if err := ValidateSubjectDescription(req.Description); err != nil {
		return err
	}

	return nil
}
