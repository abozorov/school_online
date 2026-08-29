package handler

import (
	"errors"
	"fmt"

	"github.com/abozorov/school_online/cmd/classroom/internal/models"
	"github.com/abozorov/school_online/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func responseErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, errs.ErrBadRequest),
		errors.Is(err, errs.ErrBadRequestBody),
		errors.Is(err, errs.ErrBadRequestQuery),
		errors.Is(err, models.ErrInvalidClassroomID),
		errors.Is(err, models.ErrInvalidGradeNumber),
		errors.Is(err, models.ErrInvalidClassroomLetter),
		errors.Is(err, models.ErrInvalidTeacherID),
		errors.Is(err, models.ErrInvalidAcademicYear):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("User microservice: %s", errs.ErrBadRequest))
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return status.Error(codes.DeadlineExceeded, fmt.Sprintf("User microservice: %s", errs.ErrTimeoutExceeded))
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, fmt.Sprintf("User microservice: %s", errs.ErrNotFound))
	default:
		return status.Error(codes.Internal, fmt.Sprintf("User microservice: %s", errs.ErrSomethingWentWrong))
	}
}
