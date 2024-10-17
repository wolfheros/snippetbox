package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Any code here will execute on the way down the chain
		next.ServeHTTP(w, r)
		// Any code here will execute on the way back up the chain fron handler.
	})
}

// log out all the request detail, the reason why here, attached this method to `application`,
// its trying to use [main.infoLog][io.EOF] method.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recover from panic funtion
// About Recover from Panic(pattern):
// 1. Must register a funton as defer funtion to handle a potential panic
// 2. In `defer` this funtion, call buildin function recover() within an `if` statement check `no-nil` value was found
// 3. recover() function can NOT handle children goroutine.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// handle potential panic from `next.ServerHTTP(w,r)`
		defer func() {

			// recover method, processing
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
