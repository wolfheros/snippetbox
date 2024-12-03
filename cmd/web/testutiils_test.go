package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/snippetbox/pkg/models/mock"
)

// Define a regular expression which captures the CSRF token value from the HTML for our user sigup page
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	// extract the csrf token from html body, the return is a array with match result in the first position.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

// creating a application test helper, which return a instance of `application` struct
func newTestApplication(t *testing.T) *application {
	//Create an instance of the temple cache
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance, with the same settings as productions
	session := sessions.New([]byte("3dsm5Mnygddddddddddddddddddddddd"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
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

// creat a post form method for sending POST requests to the test server.
// the final parameter to this method is aurl.Values object which can contain
// any data that you want to send in the request body.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body

}
