package handywares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
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

	errCat := bricks.ErrUnknown
	if errIsCanceled(err) {
		errCat = bricks.ErrCanceled
	} else if errIsTimeout(err) {
		errCat = bricks.ErrDeadlineExceeded
	}

	if berr, ok := httpStatusToBricksErr[code]; ok {
		errCat = berr
	}

	return rsp, errors.Join(errCat, err, fmt.Errorf("http status (%d): %s", code, http.StatusText(code)))
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
