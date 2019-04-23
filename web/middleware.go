package web

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// HTTP request logging logic

type logContextKey struct{}

var logKey = logContextKey{}

type requestInfo struct {
	id       string
	protocol string
	method   string
	scheme   string
	host     string
	uri      string
	remote   string
}

// Recoverer returns a middleware recovery handler using our standard logger
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				var err error
				switch x := rvr.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown error")
				}

				log.Error().Stack().Err(errors.Wrap(err, "unexpected error")).Msg("unexpected error")
				debug.PrintStack()
				render.Render(w, r, ErrServerError(err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Logger returns a middleware logger handler using our standard logger
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := newLogInfo(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()

		defer func() {
			info.Write(ww.Status(), ww.BytesWritten(), time.Since(t1))
		}()

		next.ServeHTTP(ww, withLogInfo(r, info))
	})
}

func getLogInfo(r *http.Request) requestInfo {
	info, _ := r.Context().Value(logKey).(requestInfo)
	return info
}

func withLogInfo(r *http.Request, info requestInfo) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), logKey, info))
}

func newLogInfo(r *http.Request) requestInfo {
	info := requestInfo{id: middleware.GetReqID(r.Context()), protocol: r.Proto, method: r.Method, scheme: "http",
		host: r.Host, uri: r.RequestURI, remote: r.RemoteAddr}
	if r.TLS != nil {
		info.scheme = "https"
	}
	return info
}

func (r *requestInfo) Write(status, bytes int, elapsed time.Duration) {
	event := log.Info()
	if status >= 500 {
		event = log.Error()
	} else if status >= 400 {
		event = log.Warn()
	}

	event.Str("id", r.id).Str("proto", r.protocol).Str("method", r.method).Str("scheme", r.scheme).Str("host", r.host).Str("remote", r.remote).Int("status", status).Int("bytes", bytes).Dur("elapsed", elapsed).Msg(r.uri)
}
