package services

import (
	"fmt"

	"github.com/abozorov/school_online/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCToErrs(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrSomethingWentWrong, err.Error())
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %s", errs.ErrBadRequestBody, st.Message())
	case codes.NotFound:
		return fmt.Errorf("%w: %s", errs.ErrNotFound, st.Message())
	case codes.AlreadyExists:
		return fmt.Errorf("%w: %s", errs.ErrAlreadyExists, st.Message())
	case codes.DeadlineExceeded, codes.Unavailable, codes.Canceled:
		return fmt.Errorf("%w: %s", errs.ErrTimeoutExceeded, st.Message())
	case codes.PermissionDenied:
		return fmt.Errorf("%w: %s", errs.ErrBadRequest, st.Message())
	default:
		return fmt.Errorf("%w: %s", errs.ErrSomethingWentWrong, st.Message())
	}
}
