package http_server

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// GetRecoverer returns a http middleware catching panics in http.Handler.
func GetRecoverer(log logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if log != nil {
						log.WithField("panic", true).Errorf("%s \n %s", rvr, string(debug.Stack()))
					}
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
