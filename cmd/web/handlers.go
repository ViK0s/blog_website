package main

import (
	"context"
	"net/http"

	base_html "blog.michalg.net/ui/html"
	"blog.michalg.net/ui/html/pages"
	"blog.michalg.net/ui/html/partials"
)

func (app *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	component := base_html.HTML("Home", pages.Hometmp(app.blogpost.Latest()), partials.Nav())
	component.Render(context.Background(), w)
}

func about(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	component := base_html.HTML("About", pages.Abouttmp(), partials.Nav())
	component.Render(context.Background(), w)

}
