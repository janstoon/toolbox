package handywares_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/handywares"
)

func TestHttpMiddlewareStackFunctionality(t *testing.T) {
	var mws handywares.HttpMiddlewareStack
	mws.
		Push(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				_, _ = rw.Write([]byte("top middleware."))

				next.ServeHTTP(rw, req)
			})
		}).
		Push(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				_, _ = rw.Write([]byte("middle middleware."))

				next.ServeHTTP(rw, req)
			})
		}).
		Push(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				_, _ = rw.Write([]byte("bottom middleware."))

				next.ServeHTTP(rw, req)
			})
		})

	srv := httptest.NewServer(mws(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("actual handler."))
	})))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	rsp, err := srv.Client().Do(req)
	srv.Close()

	require.NoError(t, err)
	assert.NotEmpty(t, rsp)
	assert.Equal(t, http.StatusOK, tricks.PtrVal(rsp).StatusCode)
	bb, err := io.ReadAll(tricks.PtrVal(rsp).Body)
	require.NoError(t, err)
	assert.NotEmpty(t, bb)
	require.NoError(t, rsp.Body.Close())
	assert.Equal(t, "top middleware.middle middleware.bottom middleware.actual handler.", string(bb))
}

func TestHttpMiddlewarePanicRecover(t *testing.T) {
	var mws handywares.HttpMiddlewareStack
	mws.PushPanicRecover()

	srv := httptest.NewServer(mws(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		panic("server panic")
	})))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	rsp, err := srv.Client().Do(req)
	srv.Close()

	require.NoError(t, err)
	assert.NotEmpty(t, rsp)
	assert.Equal(t, http.StatusInternalServerError, tricks.PtrVal(rsp).StatusCode)
	require.NoError(t, rsp.Body.Close())
}
