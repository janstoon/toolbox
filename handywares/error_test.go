package handywares_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/handywares"
)

func TestHttpStatusToBricksError(t *testing.T) {
	assert.NoError(t, handywares.HttpStatusToBricksError(http.StatusOK, nil))

	err := handywares.HttpStatusToBricksError(http.StatusNotFound, errors.New("missing"))
	require.ErrorIs(t, err, bricks.ErrNotFound)
	assert.ErrorContains(t, err, "missing")

	err = handywares.HttpStatusToBricksError(http.StatusTeapot, errors.New("odd"))
	require.ErrorIs(t, err, bricks.ErrUnknown)

	rsp, err := handywares.HttpToBricksErrorMapper(&http.Response{StatusCode: http.StatusConflict}, nil)
	require.ErrorIs(t, err, bricks.ErrAlreadyExists)
	assert.Equal(t, http.StatusConflict, rsp.StatusCode)
}

func TestBricksErrorToHttpStatusMapper(t *testing.T) {
	assert.Equal(t, http.StatusForbidden, handywares.BricksErrorToHttpStatusMapper(bricks.ErrPermissionDenied))
	assert.Equal(t, http.StatusInternalServerError, handywares.BricksErrorToHttpStatusMapper(errors.New("plain")))
	assert.Equal(t, http.StatusInternalServerError, handywares.BricksErrorToHttpStatusMapper(bricks.ErrInternal))
}

func TestGrpcErrorMappers(t *testing.T) {
	assert.NoError(t, handywares.BricksToGrpcErrorMapper(nil))
	err := handywares.BricksToGrpcErrorMapper(bricks.ErrNotFound)
	assert.Equal(t, codes.NotFound, status.Code(err))

	err = handywares.BricksToGrpcErrorMapper(errors.New("plain"))
	assert.Equal(t, codes.Unknown, status.Code(err))

	assert.NoError(t, handywares.GrpcToBricksErrorMapper(nil))
	mapped := handywares.GrpcToBricksErrorMapper(status.Error(codes.Unavailable, "down"))
	require.ErrorIs(t, mapped, bricks.ErrUnavailable)
}

func TestRedisToBricksError(t *testing.T) {
	assert.NoError(t, handywares.RedisToBricksError(nil))
	require.ErrorIs(t, handywares.RedisToBricksError(redis.Nil), bricks.ErrNotFound)
	require.ErrorIs(t, handywares.RedisToBricksError(redis.ErrClosed), bricks.ErrUnavailable)
	require.ErrorIs(t, handywares.RedisToBricksError(redis.TxFailedErr), bricks.ErrInternal)
	require.ErrorIs(t, handywares.RedisToBricksError(context.Canceled), bricks.ErrCanceled)
}

func ExampleHttpStatusToBricksError() {
	err := handywares.HttpStatusToBricksError(http.StatusNotFound, nil)
	fmt.Println(errors.Is(err, bricks.ErrNotFound))
	fmt.Println(handywares.BricksErrorToHttpStatusMapper(bricks.ErrNotFound))

	// Output:
	// true
	// 404
}
