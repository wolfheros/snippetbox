package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/snippetbox/pkg/forms"
	"github.com/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	///// No need, Pat match "/" exactly
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// panic("oops!  something went wrong") // testing

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// new cache method
	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat can't handle colon from the named capture, so need `:id` get value from query.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		switch err {
		case models.ErrNoRecord:
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}


	// new cache method
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// new empty forms.Form object to the template
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// processing the parse form, support POST, PUT, PATCH
	// anything error happen, send 400 bad request to client
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// check every field in the form satisfy the requrement. 
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// if any satisfy is missed, there will be errors in form instance.
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// create a new snippet record in the database using form data.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// put() add string value as session data
	app.session.Put(r, "flash", "Snippet sucessfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Display the user signup form")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request){
	fmt. Fprintln(w, "Create a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Display the user login form...")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "AUthenticate and login the user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Logout the user...")
}