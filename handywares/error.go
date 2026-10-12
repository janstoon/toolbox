package handywares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var httpStatusToBricksErr = map[int]error{
	http.StatusBadRequest:                   bricks.ErrInvalidArgument,
	http.StatusUnauthorized:                 bricks.ErrUnauthenticated,
	http.StatusForbidden:                    bricks.ErrPermissionDenied,
	http.StatusNotFound:                     bricks.ErrNotFound,
	http.StatusRequestTimeout:               bricks.ErrDeadlineExceeded,
	http.StatusConflict:                     bricks.ErrAlreadyExists,
	http.StatusPreconditionFailed:           bricks.ErrFailedPrecondition,
	http.StatusRequestedRangeNotSatisfiable: bricks.ErrOutOfRange,
	http.StatusTooManyRequests:              bricks.ErrResourceExhausted,

	http.StatusInternalServerError: bricks.ErrInternal,
	http.StatusNotImplemented:      bricks.ErrUnimplemented,
	http.StatusServiceUnavailable:  bricks.ErrUnavailable,
	http.StatusGatewayTimeout:      bricks.ErrDeadlineExceeded,
	http.StatusInsufficientStorage: bricks.ErrResourceExhausted,
}

// HttpStatusToBricksError maps an HTTP status to a bricks error and joins it with err.
// A nil err and a status below 400 return nil. An unmapped status joins bricks.ErrUnknown
// when err does not already carry a category from bricks.ParseError.
func HttpStatusToBricksError(code int, err error) error {
	if err == nil && code < http.StatusBadRequest {
		return nil
	}

	ferr := bricks.ParseError(err, bricks.ErrUnknown)
	if errByCode, ok := httpStatusToBricksErr[code]; ok {
		ferr = errors.Join(errByCode, err)
	}

	return errors.Join(ferr, fmt.Errorf("http status (%d): %s", code, http.StatusText(code)))
}

// HttpToBricksErrorMapper is an HTTP tripperware mapper. It returns rsp unchanged and
// replaces err with the bricks error for rsp.StatusCode.
func HttpToBricksErrorMapper(rsp *http.Response, err error) (*http.Response, error) {
	return rsp, HttpStatusToBricksError(tricks.PtrVal(rsp).StatusCode, err)
}

var bricksErrCodeToHttpStatus = map[int]int{
	bricks.ErrCodeInvalidArgument:    http.StatusBadRequest,
	bricks.ErrCodeUnauthenticated:    http.StatusUnauthorized,
	bricks.ErrCodePermissionDenied:   http.StatusForbidden,
	bricks.ErrCodeNotFound:           http.StatusNotFound,
	bricks.ErrCodeAlreadyExists:      http.StatusConflict,
	bricks.ErrCodeFailedPrecondition: http.StatusPreconditionFailed,
	bricks.ErrCodeOutOfRange:         http.StatusRequestedRangeNotSatisfiable,

	bricks.ErrCodeUnimplemented: http.StatusNotImplemented,
	bricks.ErrCodeUnavailable:   http.StatusServiceUnavailable,
}

// BricksErrorToHttpStatusMapper returns the HTTP status for a bricks error.
// An error without a known code maps to 500.
func BricksErrorToHttpStatusMapper(err error) int {
	code := http.StatusInternalServerError
	var coded bricks.Coded
	if errors.As(err, &coded) {
		if v, ok := bricksErrCodeToHttpStatus[coded.Code()]; ok {
			code = v
		}
	}

	return code
}

// BricksToGrpcErrorMapper returns a gRPC status error whose code is the bricks code.
// A nil err returns nil. An error without a code uses codes.Unknown.
func BricksToGrpcErrorMapper(err error) error {
	if err == nil {
		return nil
	}

	code := codes.Unknown
	var coded bricks.Coded
	if errors.As(err, &coded) {
		code = codes.Code(coded.Code())
	}

	return status.Error(code, err.Error())
}

var grpcCodeToBricksErr = map[codes.Code]error{
	codes.Canceled:           bricks.ErrCanceled,
	codes.Unknown:            bricks.ErrUnknown,
	codes.InvalidArgument:    bricks.ErrInvalidArgument,
	codes.DeadlineExceeded:   bricks.ErrDeadlineExceeded,
	codes.NotFound:           bricks.ErrNotFound,
	codes.AlreadyExists:      bricks.ErrAlreadyExists,
	codes.PermissionDenied:   bricks.ErrPermissionDenied,
	codes.ResourceExhausted:  bricks.ErrResourceExhausted,
	codes.FailedPrecondition: bricks.ErrFailedPrecondition,
	codes.Aborted:            bricks.ErrAborted,
	codes.OutOfRange:         bricks.ErrOutOfRange,
	codes.Unimplemented:      bricks.ErrUnimplemented,
	codes.Internal:           bricks.ErrInternal,
	codes.Unavailable:        bricks.ErrUnavailable,
	codes.DataLoss:           bricks.ErrDataLoss,
	codes.Unauthenticated:    bricks.ErrUnauthenticated,
}

// GrpcToBricksErrorMapper joins the bricks error for the gRPC status code with err.
// A nil err returns nil.
func GrpcToBricksErrorMapper(err error) error {
	if err == nil {
		return nil
	}

	berr := bricks.ErrUnknown
	code := status.Code(err)
	if v, ok := grpcCodeToBricksErr[code]; ok {
		berr = v
	}

	return errors.Join(berr, err)
}

// RedisToBricksError maps redis.Nil to ErrNotFound, redis.ErrClosed to ErrUnavailable,
// and redis.TxFailedErr to ErrInternal. Other errors go through bricks.ParseError.
// A nil err returns nil.
func RedisToBricksError(err error) error {
	if err == nil {
		return nil
	}

	errW := bricks.ErrUnknown
	if errors.Is(err, redis.Nil) {
		errW = errors.Join(bricks.ErrNotFound, err)
	} else if errors.Is(err, redis.ErrClosed) {
		errW = errors.Join(bricks.ErrUnavailable, err)
	} else if errors.Is(err, redis.TxFailedErr) {
		errW = errors.Join(bricks.ErrInternal, err)
	} else {
		errW = bricks.ParseError(err, bricks.ErrUnknown)
	}

	return errW
}
