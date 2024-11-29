package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// creating a application test helper, which return a instance of `application` struct
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
	}

}

// Define a custom http test server struct
// !!!Embedded Field!!!
// Be careful here, type struct have noname field (*httptest.Server). which mean it can be directily use.
type testServer struct {
	*httptest.Server
}

// Start test server, using app.routes as the handler routes
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	//initialize a new cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	//add cookie jar to the client,
	//so that response cookies are stored
	// then sent within subsequent requests
	ts.Client().Jar = jar

	//Disable rediect-following for the client.
	// instead we want test first https response sent by our server.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

// Implement a get method on our custom testServer type.
// This make a GET request
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	// here make a GET request to testServer
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
