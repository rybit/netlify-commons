package tracing

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type contextKey string

const tracerKey = contextKey("nf-tracer-key")

func WrapWithTracer(r *http.Request, rt *RequestTracer) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), tracerKey, rt))
}

func GetFromContext(ctx context.Context) *RequestTracer {
	val := ctx.Value(tracerKey)
	if val == nil {
		return nil
	}
	entry, ok := val.(*RequestTracer)
	if ok {
		return entry
	}
	return nil
}

func GetTracer(r *http.Request) *RequestTracer {
	return GetFromContext(r.Context())
}

func GetRequestID(r *http.Request) string {
	if id := GetRequestIDFromContext(r.Context()); id != "" {
		return id
	}
	if rid := r.Header.Get(HeaderRequestUUID); rid != "" {
		return rid
	}
	return uuid.NewV4().String()
}

func GetRequestIDFromContext(ctx context.Context) string {
	tr := GetFromContext(ctx)
	if tr == nil {
		return ""
	}
	return tr.RequestID
}
