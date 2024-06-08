package handywares

import (
	"errors"
	"fmt"
	"net/http"

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

	errCat := bricks.ParseError(err, bricks.ErrUnknown)
	if berr, ok := httpStatusToBricksErr[code]; ok {
		errCat = errors.Join(berr, err)
	}

	return rsp, errors.Join(errCat, fmt.Errorf("http status (%d): %s", code, http.StatusText(code)))
}
