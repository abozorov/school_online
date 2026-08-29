package handler

import (
	"errors"
	"fmt"

	"github.com/abozorov/school_online/cmd/raiting/internal/models"
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
		errors.Is(err, models.ErrInvalidJournalID):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("Raiting microservice: %s", errs.ErrBadRequest))
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return status.Error(codes.DeadlineExceeded, fmt.Sprintf("Raiting microservice: %s", errs.ErrTimeoutExceeded))
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, fmt.Sprintf("Raiting microservice: %s", errs.ErrNotFound))
	default:
		return status.Error(codes.Internal, fmt.Sprintf("Raiting microservice: %s", errs.ErrSomethingWentWrong))
	}
}
