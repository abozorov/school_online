package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/abozorov/school_online/pkg/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func postgresToErrs(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("error %w: %w", err, errs.ErrNotFound)
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)
		case "23503":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)
		case "23502":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)
		case "23514":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)
		case "22001":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)
		case "22003":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)
		case "57014":
			return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)
		case "22P02":
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)
		}
	}

	return fmt.Errorf("error %w: %w", err, errs.ErrSomethingWentWrong)
}
