package handywares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
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
