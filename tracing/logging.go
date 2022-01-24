package tracing

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func RequestLogger(r *http.Request, log logrus.FieldLogger) (logrus.FieldLogger, string) {
	if r.Header.Get(HeaderNFDebugLogging) != "" {
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		if entry, ok := log.(*logrus.Entry); ok {
			log = logger.WithFields(entry.Data)
			logger.Hooks = entry.Logger.Hooks
		}
	}

	reqID := r.Header.Get(HeaderRequestUUID)
	if reqID == "" {
		reqID = uuid.NewV4().String()
		r.Header.Set(HeaderRequestUUID, reqID)
	}

	log = log.WithFields(logrus.Fields{
		"request_id": reqID,
	})

	return log, reqID
}

func GetLoggerFromContext(ctx context.Context) logrus.FieldLogger {
	entry := GetFromContext(ctx)
	if entry == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return entry.FieldLogger
}

func GetLogger(r *http.Request) logrus.FieldLogger {
	return GetLoggerFromContext(r.Context())
}

// SetLogField will add the field to this log line and every one following
func SetLogField(r *http.Request, key string, value interface{}) logrus.FieldLogger {
	entry := GetTracer(r)
	if entry == nil {
		return logrus.StandardLogger().WithField(key, value)
	}
	return entry.SetLogField(key, value)
}

// SetLogFields will add the fields to this log line and every one following
func SetLogFields(r *http.Request, fields logrus.Fields) logrus.FieldLogger {
	entry := GetTracer(r)
	if entry == nil {
		return logrus.StandardLogger().WithFields(fields)
	}

	return entry.SetLogFields(fields)
}

// SetFinalField will add a field to the canonical line created at in Finish. It will add
// it to this line, but not every log line in between
func SetFinalField(r *http.Request, key string, value interface{}) logrus.FieldLogger {
	return SetFinalFieldWithContext(r.Context(), key, value)
}

// SetFinalFieldWithContext will add a field to the canonical line created at in Finish. It will add
// it to this line, but not every log line in between
func SetFinalFieldWithContext(ctx context.Context, key string, value interface{}) logrus.FieldLogger {
	entry := GetFromContext(ctx)
	if entry == nil {
		return logrus.StandardLogger().WithField(key, value)
	}
	return entry.SetFinalField(key, value)
}
