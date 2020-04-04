package routes

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Prevent abnormal shutdown while panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Println(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
