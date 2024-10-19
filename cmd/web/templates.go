package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/snippetbox/pkg/forms"
	"github.com/snippetbox/pkg/models"
) // internal package - domainname/projectname/filepath/packagename

// use it hold more than one data during parse template, cause template only
// support one data in each template.
type templateData struct {
	// current year
	CurrentYear int

	Form *forms.Form

	Snippet *models.Snippet

	Snippets []*models.Snippet
}

// cache func, using `map[string]*template.Template` to cache parsed template from files.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	//
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// return the last element of the filepath, which is the filename
		name := filepath.Base(page)

		// create a template set then add functions to it for invoke in the template.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add any `layout` templates to `page` layout,which is parsed before.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// add `partial` templates to `page` layout too
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// add the templates to the cache map, using `layout` name, such as `home.layout.tmpl` as the key.
		cache[name] = ts

	}
	return cache, nil
}

// format time for template use
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:03")
}

// what it did in here, its trying to invoke the function in the *.tmpl file to implement dynamic data.
var functions = template.FuncMap{
	"humanDate": humanDate,
}
