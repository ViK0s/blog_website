package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	base_html "blog.michalg.net/ui/html"
	"blog.michalg.net/ui/html/pages"
	"blog.michalg.net/ui/html/partials"
)

func (app *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// jesus
	component := base_html.HTML("Home", pages.Hometmp(partials.LatestBlogPostTempl(app.blogpost.Latest()[0]), partials.Another_Section(app.blogpost.Latest()[1])), partials.Nav())
	component.Render(context.Background(), w)
}

func (app *app) about(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	component := base_html.HTML("About", pages.Abouttmp(partials.Contact_Section()), partials.Nav())
	component.Render(context.Background(), w)

}

func (app *app) blogpostview(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		fmt.Println(id)
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.blogpost.Get(id)

	if err != nil || id < 1 {
		fmt.Println(err)
		return
	}

	// if err != nil {
	//     if errors.Is(err, models.ErrNoRecord) {
	//         app.notFound(w)
	//     } else {
	//         app.serverError(w, err)
	//     }
	//     return
	// }
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Println(id)
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *app) projectshandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/projects" {
		http.NotFound(w, r)
		return
	}

	component := base_html.HTML("Projects", pages.ProjectsSite(partials.Projecttemp(app.project.Latest())), partials.Nav())
	component.Render(context.Background(), w)
}

func (app *app) blogpostpage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.NotFound(w, r)
		return
	}

	component := base_html.HTML("Blogs", pages.Blogsite(partials.BlogPostFull(app.blogpost.Latest()), partials.Archive(app.blogpost.Latest())), partials.Nav())
	component.Render(context.Background(), w)
}
