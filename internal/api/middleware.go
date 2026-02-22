package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware Logger ----
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (app *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rec, r)
		dur := time.Since(start)
		log.Printf("%s | %d | %s | %s", r.Method, rec.statusCode, formatDuration(dur), r.URL.Path)
	})
}

func formatDuration(d time.Duration) string {
	var s string
	if d >= time.Second {
		s = fmt.Sprintf("%ds", int64(d.Seconds()+0.5))
	} else if d >= time.Millisecond {
		ms := int64(float64(d.Nanoseconds())/1e6 + 0.5)
		s = fmt.Sprintf("%dms", ms)
	} else if d >= time.Microsecond {
		us := int64(float64(d.Nanoseconds())/1e3 + 0.5)
		s = fmt.Sprintf("%dµs", us)
	} else {
		s = fmt.Sprintf("%dns", d.Nanoseconds())
	}
	return fmt.Sprintf("%6s", s)
}

//----
