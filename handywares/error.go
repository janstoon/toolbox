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

func HttpToBricksErrorMapper(rsp *http.Response, err error) (*http.Response, error) {
	code := tricks.PtrVal(rsp).StatusCode
	if err == nil && code < http.StatusBadRequest {
		return rsp, nil
	}

	errCat := bricks.ParseError(err, bricks.ErrUnknown)
	if berr, ok := httpStatusToBricksErr[code]; ok {
		errCat = errors.Join(berr, err)
	}

	return rsp, errors.Join(errCat, fmt.Errorf("http status (%d): %s", code, http.StatusText(code)))
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
