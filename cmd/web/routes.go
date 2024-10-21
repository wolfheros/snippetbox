package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// create a middleware chain containing our `standard` minddlewares.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// this middleware is for session management
	dynamicMiddleware := alice.New(app.session.Enable)
	// ServerMux also implemented [Handler]interface.
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// !!!WARNING!!!
	// Pat match patterns in the order that they are registerd, and `/snippet/create` equal to `/snippet/:id`,
	// that's why `/snippet/create` need registed first before `/snippet/:id`, otherwise it will become part of pattern
	// `/snippet/:id`
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// static pattern no need change against third-party router Pat.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
