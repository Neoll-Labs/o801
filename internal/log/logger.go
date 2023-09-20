/*
 license x
*/

package log

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%5s %20s %10s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
