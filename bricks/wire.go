package bricks

import (
	"context"
	"errors"
	"os"
)

const (
	// ErrCodeCanceled is the status code for a canceled operation.
	ErrCodeCanceled = iota + 1
	// ErrCodeUnknown is the status code for an error that has no more specific code.
	ErrCodeUnknown
	// ErrCodeInvalidArgument is the status code for a rejected argument.
	ErrCodeInvalidArgument
	// ErrCodeDeadlineExceeded is the status code for an operation that timed out.
	ErrCodeDeadlineExceeded
	// ErrCodeNotFound is the status code for a missing entity.
	ErrCodeNotFound
	// ErrCodeAlreadyExists is the status code for an entity that is already present.
	ErrCodeAlreadyExists
	// ErrCodePermissionDenied is the status code for an identified caller who lacks permission.
	ErrCodePermissionDenied
	// ErrCodeResourceExhausted is the status code for a quota or capacity failure.
	ErrCodeResourceExhausted
	// ErrCodeFailedPrecondition is the status code for a call made in the wrong state.
	ErrCodeFailedPrecondition
	// ErrCodeAborted is the status code for an operation aborted by a conflict.
	ErrCodeAborted
	// ErrCodeOutOfRange is the status code for a value past the valid range.
	ErrCodeOutOfRange
	// ErrCodeUnimplemented is the status code for an unsupported operation.
	ErrCodeUnimplemented
	// ErrCodeInternal is the status code for a broken invariant in the service.
	ErrCodeInternal
	// ErrCodeUnavailable is the status code for a service that cannot be reached.
	ErrCodeUnavailable
	// ErrCodeDataLoss is the status code for unrecoverable data loss.
	ErrCodeDataLoss
	// ErrCodeUnauthenticated is the status code for a caller who was not identified.
	ErrCodeUnauthenticated
)

var (
	// ErrRetryable marks an error that the caller may try again.
	ErrRetryable = errors.New("temporary issue")

	// ErrUnknown is the fallback error when no more specific category applies.
	ErrUnknown = ErrorWithCode(ErrCodeUnknown, errors.New("unknown error"))
	// ErrCanceled reports that the operation was canceled. It is retryable.
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

	// ErrCustomerSide groups errors caused by the caller.
	ErrCustomerSide = errors.New("customer-side error")
	// ErrInvalidArgument reports a rejected argument.
	ErrInvalidArgument = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeInvalidArgument, errors.New("invalid argument")))
	// ErrUnauthenticated reports that the caller was not identified.
	ErrUnauthenticated = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeUnauthenticated, errors.New("caller not identified (unauthenticated)")))
	// ErrPermissionDenied reports that an identified caller lacks permission.
	ErrPermissionDenied = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodePermissionDenied, errors.New("caller identified but permission denied")))
	// ErrNotFound reports that the requested entity does not exist.
	ErrNotFound = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeNotFound, errors.New("requested entity was not found")))
	// ErrAlreadyExists reports that the entity is already present.
	ErrAlreadyExists = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeAlreadyExists, errors.New("entity already exists")))
	// ErrFailedPrecondition reports that the system is not in the state the operation requires.
	ErrFailedPrecondition = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeFailedPrecondition, errors.New("failed precondition")))
	// ErrOutOfRange reports that a value was past the valid range.
	ErrOutOfRange = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeOutOfRange, errors.New("operation was attempted past the valid range")))
	// ErrAborted reports that the operation was aborted.
	ErrAborted = errors.Join(ErrCustomerSide, ErrorWithCode(ErrCodeAborted, errors.New("operation was aborted")))
	// ErrDataLoss reports unrecoverable data loss or corruption.
	ErrDataLoss = errors.Join(ErrCustomerSide,
		ErrorWithCode(ErrCodeDataLoss, errors.New("unrecoverable data loss or corruption")))

	// ----------------------
	//     Supplier-side
	// ----------------------

	// ErrSupplierSide groups errors caused by the service.
	ErrSupplierSide = errors.New("supplier-side error")
	// ErrInternal reports a broken invariant in the service. It is retryable.
	ErrInternal = errors.Join(ErrRetryable, ErrSupplierSide,
		ErrorWithCode(ErrCodeInternal, errors.New("some invariants expected by underlying system has been broken")))
	// ErrUnimplemented reports an operation the service does not support.
	ErrUnimplemented = errors.Join(ErrSupplierSide,
		ErrorWithCode(ErrCodeUnimplemented,
			errors.New("operation is not implemented or not supported/enabled in this service")))
	// ErrUnavailable reports that the service cannot be reached. It is retryable.
	ErrUnavailable = errors.Join(ErrRetryable, ErrSupplierSide,
		ErrorWithCode(ErrCodeUnavailable, errors.New("service is currently unavailable")))
)

// Coded is an error that carries an ErrCode value.
type Coded interface {
	// Code returns the ErrCode stored on the error.
	Code() int
}

// ErrorWithCode returns err annotated with code. errors.As can recover the Coded value.
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

// Code returns the status code stored by ErrorWithCode.
func (ce codedError) Code() int {
	return ce.code
}

// Unwrap returns the error passed to ErrorWithCode.
func (ce codedError) Unwrap() error {
	return ce.error
}

// ParseError classifies err as canceled, timed out, or unknown.
// A nil err returns nil. unknown is used when err is neither a cancellation nor a timeout.
// The returned error joins that category with err.
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

// MessageEnvelope is a portable wrapper for a message body and its delivery metadata.
type MessageEnvelope struct {
	// Id identifies the message.
	Id string
	// Retried is how many times delivery has already been attempted.
	Retried uint
	// Note is the message body.
	Note []byte
}
