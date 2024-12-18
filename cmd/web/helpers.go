package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
	//"github.com/snippetbox/pkg/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// render the template, first it will check the cache first, then it will
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	//
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// initial a new buffer for save the response in case error happened.
	buf := new(bytes.Buffer)

	// write the template to the buffer, instead of straight to the http.ResponseWritter.
	// if there is a error in server, call our serverError helper and then return.
	// err := ts.Execute(buf, td)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	err := ts.Execute(w, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

// add default data, current year, authentirate or not to templateData,
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	// return app.session.Exists(r, "authenticatedUserID")

	// instead of check the session data, here check the context from request whether the user
	// is authenticated, the request context has been change from authenticated middleware through routes.
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
