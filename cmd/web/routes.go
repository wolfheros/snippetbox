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
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
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
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Authencate routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	//Add a new testing GET /ping route
	mux.Get("/ping", http.HandlerFunc(ping))

	// static pattern no need change against third-party router Pat.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
