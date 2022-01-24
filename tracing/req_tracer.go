package tracing

import (
	"net/http"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

type RequestTracer struct {
	*trackingWriter
	logrus.FieldLogger

	RequestID   string
	finalFields map[string]interface{}

	remoteAddr  string
	method      string
	originalURL string
	referrer    string
	span        opentracing.Span
	start       time.Time
}

func NewTracer(w http.ResponseWriter, r *http.Request, log logrus.FieldLogger, service, resource string) (http.ResponseWriter, *http.Request, *RequestTracer) {
	var reqID string
	log, reqID = RequestLogger(r, log)

	r, span := WrapWithSpan(r, reqID, service, resource)
	trackWriter := &trackingWriter{
		writer: w,
		log:    log,
	}

	rt := &RequestTracer{
		originalURL: r.URL.String(),
		method:      r.Method,
		referrer:    r.Referer(),
		remoteAddr:  r.RemoteAddr,

		RequestID:      reqID,
		span:           span,
		trackingWriter: trackWriter,
		FieldLogger:    log,
		finalFields:    make(map[string]interface{}),
	}
	r = WrapWithTracer(r, rt)

	return rt, r, rt
}

func (rt *RequestTracer) Start() {
	rt.start = time.Now()
	rt.WithFields(logrus.Fields{
		"method":      rt.method,
		"remote_addr": rt.remoteAddr,
		"referer":     rt.referrer,
		"url":         rt.originalURL,
	}).Info("Starting Request")
}

func (rt *RequestTracer) Finish() {
	dur := time.Since(rt.start)

	fields := logrus.Fields{}
	for k, v := range rt.finalFields {
		fields[k] = v
	}

	fields["status_code"] = rt.trackingWriter.status
	fields["rsp_bytes"] = rt.trackingWriter.rspBytes
	fields["url"] = rt.originalURL
	fields["method"] = rt.method
	fields["dur"] = dur.String()
	fields["dur_ns"] = dur.Nanoseconds()

	// Setting the status as an int doesn't propogate for use in datadog dashboards,
	// so we convert to a string.
	rt.span.SetTag(ext.HTTPCode, strconv.Itoa(rt.trackingWriter.status))
	rt.span.Finish()
	rt.WithFields(fields).Info("Completed Request")
}

func (rt *RequestTracer) SetLogField(key string, value interface{}) logrus.FieldLogger {
	rt.FieldLogger = rt.FieldLogger.WithField(key, value)
	return rt.FieldLogger
}

func (rt *RequestTracer) SetLogFields(fields logrus.Fields) logrus.FieldLogger {
	rt.FieldLogger = rt.FieldLogger.WithFields(fields)
	return rt.FieldLogger
}

func (rt *RequestTracer) SetFinalField(key string, value interface{}) logrus.FieldLogger {
	rt.finalFields[key] = value
	return rt.FieldLogger.WithField(key, value)
}
