package handywares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"

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

	rsp, err := srv.Client().Get(srv.URL)
	srv.Close()

	assert.NoError(t, err)
	assert.NotEmpty(t, rsp)
	assert.Equal(t, http.StatusOK, tricks.PtrVal(rsp).StatusCode)
	bb, err := io.ReadAll(tricks.PtrVal(rsp).Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, bb)
	assert.Equal(t, string(bb), "top middleware.middle middleware.bottom middleware.actual handler.")
}

func TestHttpMiddlewarePanicRecover(t *testing.T) {
	var mws handywares.HttpMiddlewareStack
	mws.PushPanicRecover()

	srv := httptest.NewServer(mws(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		panic("server panic")
	})))

	rsp, err := srv.Client().Get(srv.URL)
	srv.Close()

	assert.NoError(t, err)
	assert.NotEmpty(t, rsp)
	assert.Equal(t, http.StatusInternalServerError, tricks.PtrVal(rsp).StatusCode)
}
