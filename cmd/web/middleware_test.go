package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeader(t *testing.T) {
	// initialize a new httptest.ResponseRecords and dummy http.Request
	rr := httptest.NewRecorder()
	// dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a mock HTTP handler for testing
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	//
	secureHeaders(next).ServeHTTP(rr, r)

	//get result from http.ResponseRecorder
	rs := rr.Result()

	/**
	* Check the result
	**/

	// check the middleware has correctly set X-Frame-Options
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}
	//Check middleware has correctly set X-XSS-Protection header
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}
	//check middleware has correctly called the next handler in line
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %q; got %q", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}
