package bricks

import (
	"context"
	"errors"
	"os"
)

const (
	ErrCodeCanceled = iota + 1
	ErrCodeUnknown
	ErrCodeInvalidArgument
	ErrCodeDeadlineExceeded
	ErrCodeNotFound
	ErrCodeAlreadyExists
	ErrCodePermissionDenied
	ErrCodeResourceExhausted
	ErrCodeFailedPrecondition
	ErrCodeAborted
	ErrCodeOutOfRange
	ErrCodeUnimplemented
	ErrCodeInternal
	ErrCodeUnavailable
	ErrCodeDataLoss
	ErrCodeUnauthenticated
)

var (
	ErrRetryable = errors.New("temporary issue")

	ErrUnknown  = ErrorWithCode(ErrCodeUnknown, errors.New("unknown error"))
	ErrCanceled = errors.Join(ErrRetryable, ErrorWithCode(ErrCodeCanceled, errors.New("operation was canceled")))

	// ErrDeadlineExceeded indicates timeout has been reached regardless of operation result.
	//
	// May result in these http status codes:
	//   * 408
	//   * 504
	ErrDeadlineExceeded = errors.Join(ErrRetryable,
		ErrorWithCode(ErrCodeDeadlineExceeded, errors.New("operation expired before completion")))

	// ErrResourceExhausted indicates some resource has been exhausted,
	// perhaps a per-user quota (like rate limiting), or perhaps the entire file system is out of space.
	// Situations like out-of-memory and server overload, or when a message is larger than the configured maximum size.
	//
	// May result in these http status codes:
	//   * 429
	//   * 507
	ErrResourceExhausted = errors.Join(ErrRetryable,
		ErrorWithCode(ErrCodeResourceExhausted, errors.New("resource has been exhausted")))

	// ----------------------
	//     Customer-side
	// ----------------------

	ErrCustomerSide    = errors.New("customer-side error")
	ErrInvalidArgument = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeInvalidArgument, errors.New("invalid argument")))
	ErrUnauthenticated = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeUnauthenticated, errors.New("caller not identified (unauthenticated)")))
	ErrPermissionDenied = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodePermissionDenied, errors.New("caller identified but permission denied")))
	ErrNotFound = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeNotFound, errors.New("requested entity was not found")))
	ErrAlreadyExists = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeAlreadyExists, errors.New("entity already exists")))
	ErrFailedPrecondition = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeFailedPrecondition, errors.New("failed precondition")))
	ErrOutOfRange = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeOutOfRange, errors.New("operation was attempted past the valid range")))
	ErrAborted  = errors.Join(ErrCustomerSide, ErrorWithCode(ErrCodeAborted, errors.New("operation was aborted")))
	ErrDataLoss = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeDataLoss, errors.New("unrecoverable data loss or corruption")))

	// ----------------------
	//     Supplier-side
	// ----------------------

	ErrSupplierSide = errors.New("supplier-side error")
	ErrInternal     = errors.Join(ErrRetryable, ErrSupplierSide,
		ErrorWithCode(ErrCodeInternal, errors.New("some invariants expected by underlying system has been broken")))
	ErrUnimplemented = errors.Join(ErrSupplierSide,
		ErrorWithCode(ErrCodeUnimplemented,
			errors.New("operation is not implemented or not supported/enabled in this service")))
	ErrUnavailable = errors.Join(ErrRetryable, ErrSupplierSide,
		ErrorWithCode(ErrCodeUnavailable, errors.New("service is currently unavailable")))
)

type Coded interface {
	Code() int
}

func ErrorWithCode(code int, err error) error {
	return codedError{
		error: err,
		code:  code,
	}
}

type codedError struct {
	error
	code int
}

func (ce codedError) Code() int {
	return ce.code
}

func (ce codedError) Unwrap() error {
	return ce.error
}

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
