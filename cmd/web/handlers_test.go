package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	//initialize a new httptest.Response
	rr := httptest.NewRecorder()

	//Initialize a new dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//Call the ping handler function, passing in the httptest.ResponseRecorder
	ping(rr, r)

	/**
	*	This is how to test handlers, here are two major functionalities include:
	 */

	// <1> check that the response status code is 200
	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()

	// <2> check that the response body is "OK"
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		// can be use to fail the test,
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}
