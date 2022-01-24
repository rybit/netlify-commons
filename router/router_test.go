package router

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/netlify/netlify-commons/testutil"
	"github.com/netlify/netlify-commons/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

func TestCORS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "myexamplehost.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	t.Run("enabled", func(t *testing.T) {
		rsp := do(t, []Option{OptEnableCORS}, "", "/", nil, req)
		assert.Equal(t, http.StatusOK, rsp.Code)
	})
	t.Run("disabled", func(t *testing.T) {
		rsp := do(t, nil, "", "/", nil, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rsp.Code)
	})
}

func TestCallthrough(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	var callCount int
	handler := func(w http.ResponseWriter, r *http.Request) error {
		callCount++
		return BadRequestError("")
	}
	rsp := do(t, nil, "", "/", handler, req)
	assert.Equal(t, http.StatusBadRequest, rsp.Code)
	assert.Equal(t, 1, callCount)
}

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	scenarios := map[string]struct {
		opts []Option
		code int
	}{
		"disabled": {[]Option{OptHealthCheck("", nil)}, http.StatusNotFound},
		"default":  {[]Option{OptHealthCheck("/health", nil)}, http.StatusOK},
		"custom": {[]Option{OptHealthCheck(
			"/health",
			func(_ http.ResponseWriter, r *http.Request) error {
				return UnauthorizedError("")
			})},
			http.StatusUnauthorized},
	}

	for name, scene := range scenarios {
		t.Run(name, func(t *testing.T) {
			rsp := do(t, scene.opts, "", "/", nil, req)
			assert.Equal(t, scene.code, rsp.Code)
		})
	}
}

func TestTracing(t *testing.T) {
	og := opentracing.GlobalTracer()
	mt := mocktracer.New()
	opentracing.SetGlobalTracer(mt)
	defer func() {
		opentracing.SetGlobalTracer(og)
	}()

	noop := func(w http.ResponseWriter, r *http.Request) error {
		assert.NotNil(t, tracing.GetTracer(r))
		w.WriteHeader(http.StatusOK)
		return nil
	}

	tl, logHook := testutil.TestLogger(t)
	r := New(tl, OptEnableTracing("some-service"))

	r.Method(http.MethodPatch, "/patch", noop)
	r.Delete("/abc/{def}", noop)
	r.Get("/abc/{def}", noop)
	r.Get("/", noop)
	r.Post("/def/ghi", noop)
	r.Put("/asdf/", noop)
	r.Route("/sub", func(r Router) {
		r.Get("/path", noop)
	})

	scenes := map[string]struct {
		method, path, resourceName string
	}{
		"get":          {http.MethodGet, "/abc/def", "GET::abc.def"},
		"delete":       {http.MethodDelete, "/abc/hfj", "DELETE::abc.def"},
		"post":         {http.MethodPost, "/def/ghi", "POST::def.ghi"},
		"put":          {http.MethodPut, "/asdf/", "PUT::asdf"},
		"patch":        {http.MethodPatch, "/patch", "PATCH::patch"},
		"subroute":     {http.MethodGet, "/sub/path", "GET::sub.path"},
		"single_slash": {http.MethodGet, "/", "GET"},
	}

	for name, scene := range scenes {
		t.Run(name, func(t *testing.T) {
			mt.Reset()
			logHook.Reset()

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(scene.method, scene.path, nil))
			assert.Equal(t, http.StatusOK, rec.Code)

			spans := mt.FinishedSpans()
			if assert.Equal(t, 1, len(spans)) {
				assert.Equal(t, "some-service", spans[0].Tag(ext.ServiceName))
				assert.Equal(t, scene.resourceName, spans[0].Tag(ext.ResourceName))
				assert.Equal(t, strconv.Itoa(http.StatusOK), spans[0].Tag(ext.HTTPCode))
			}
			// should be a starting and finished request for each request
			assert.Len(t, logHook.AllEntries(), 2)
		})
	}
}

func TestVersionHeader(t *testing.T) {
	scenes := map[string]struct {
		version  string
		expected string
		header   string
		svc      string
	}{
		"custom":  {version: "123", expected: "123", header: "x-nf-something-version", svc: "something"},
		"default": {version: "", expected: "unknown", header: "x-nf-something-version", svc: "something"},
	}
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	for name, scene := range scenes {
		t.Run(name, func(t *testing.T) {
			opts := []Option{OptVersionHeader(scene.svc, scene.version)}
			rsp := do(t, opts, scene.svc, "/", nil, req)
			assert.Equal(t, scene.expected, rsp.Header().Get(scene.header), t.Name())
		})
	}
}

func do(t *testing.T, opts []Option, svcName, path string, handler APIHandler, req *http.Request) *httptest.ResponseRecorder {
	if opts == nil {
		opts = []Option{}
	}
	r := New(logrus.WithField("test", t.Name()), opts...)

	if handler == nil {
		handler = func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}
	}
	if path != "" {
		r.Get(path, handler)
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}
