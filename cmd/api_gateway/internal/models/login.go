package models

import (
	"strings"

	"github.com/abozorov/school_online/pkg/errs"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" || len(r.Password) < 8 {
		return errs.ErrBadRequestBody
	}
	return nil
}
