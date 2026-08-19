package handler

import (
	"errors"

	"github.com/abozorov/school_online/cmd/user/internal/models"
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
		errors.Is(err, models.ErrEmptyName),
		errors.Is(err, models.ErrEmptyEmail),
		errors.Is(err, models.ErrEmptyPhone),
		errors.Is(err, models.ErrEmptyUserID),
		errors.Is(err, models.ErrInvalidAge),
		errors.Is(err, models.ErrInvalidEmail),
		errors.Is(err, models.ErrInvalidName),
		errors.Is(err, models.ErrInvalidPhone),
		errors.Is(err, models.ErrInvalidUSerId),
		errors.Is(err, models.ErrInvalidID),
		errors.Is(err, models.ErrInvalidRole),
		errors.Is(err, models.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
