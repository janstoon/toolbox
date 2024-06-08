package bricks

import (
	"context"
	"errors"
	"os"
)

var (
	ErrRetryable = errors.New("temporary issue")

	ErrUnknown  = errors.New("unknown error")
	ErrCanceled = errors.Join(ErrRetryable, errors.New("operation was canceled"))

	// ErrDeadlineExceeded indicates timeout has been reached regardless of operation result.
	//
	// May result in these http status codes:
	//   * 408
	//   * 504
	ErrDeadlineExceeded = errors.Join(ErrRetryable, errors.New("operation expired before completion"))

	// ErrResourceExhausted indicates some resource has been exhausted,
	// perhaps a per-user quota (like rate limiting), or perhaps the entire file system is out of space.
	// Situations like out-of-memory and server overload, or when a message is larger than the configured maximum size.
	//
	// May result in these http status codes:
	//   * 429
	//   * 507
	ErrResourceExhausted = errors.Join(ErrRetryable, errors.New("resource has been exhausted"))

	// ----------------------
	//     Customer-side
	// ----------------------

	ErrCustomerSide       = errors.New("customer-side error")
	ErrInvalidArgument    = errors.Join(ErrCustomerSide, errors.New("invalid argument"))
	ErrUnauthenticated    = errors.Join(ErrCustomerSide, errors.New("caller not identified (unauthenticated)"))
	ErrPermissionDenied   = errors.Join(ErrCustomerSide, errors.New("caller identified but permission denied"))
	ErrNotFound           = errors.Join(ErrCustomerSide, errors.New("requested entity was not found"))
	ErrAlreadyExists      = errors.Join(ErrCustomerSide, errors.New("entity already exists"))
	ErrFailedPrecondition = errors.Join(ErrCustomerSide, errors.New("failed precondition"))
	ErrOutOfRange         = errors.Join(ErrCustomerSide, errors.New("operation was attempted past the valid range"))
	ErrAborted            = errors.Join(ErrCustomerSide, errors.New("operation was aborted"))
	ErrDataLoss           = errors.Join(ErrCustomerSide, errors.New("unrecoverable data loss or corruption"))

	// ----------------------
	//     Supplier-side
	// ----------------------

	ErrSupplierSide = errors.New("supplier-side error")
	ErrInternal     = errors.Join(ErrRetryable, ErrSupplierSide,
		errors.New("some invariants expected by underlying system has been broken"))
	ErrUnimplemented = errors.Join(ErrSupplierSide,
		errors.New("operation is not implemented or not supported/enabled in this service"))
	ErrUnavailable = errors.Join(ErrRetryable, ErrSupplierSide, errors.New("service is currently unavailable"))
)

func ParseError(err, unknown error) error {
	if err == nil {
		return nil
	}

	errCat := unknown
	if errIsCanceled(err) {
		errCat = ErrCanceled
	} else if errIsTimeout(err) {
		errCat = ErrDeadlineExceeded
	}

	return errors.Join(errCat, err)
}

func errIsTimeout(err error) bool {
	var terr interface{ Timeout() bool }

	return (errors.As(err, &terr) && terr.Timeout()) ||
		errors.Is(err, os.ErrDeadlineExceeded) ||
		errors.Is(err, context.DeadlineExceeded)
}

func errIsCanceled(err error) bool {
	return errors.Is(err, context.Canceled)
}

type MessageEnvelope struct {
	Id      string
	Retried uint
	Note    []byte
}
