package main

import (
	"context"
	"fmt"
	"log"
	"os"

	base_html "blog.michalg.net/ui/html"
	"blog.michalg.net/ui/html/pages"
	"blog.michalg.net/ui/html/partials"
)

func (app *app) build() {
	f, err := os.Create("/.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	component := base_html.HTML("Home", pages.Hometmp(partials.LatestBlogPostTempl(app.blogpost.Latest()[0]), partials.Another_Section(app.blogpost.Latest()[1])), partials.Nav())
	component.Render(context.Background(), f)

	f, err = os.Create("about.html")

	component = base_html.HTML("About", pages.Abouttmp(partials.Contact_Section()), partials.Nav())
	component.Render(context.Background(), f)

	for id := range app.blogpost.Latest()[0].ID {

		f, err = os.Create("/blog/view?id=" + "id" + ".html")
		if err != nil || id < 1 {
			fmt.Println(id)
			return
		}
		// Use the SnippetModel object's Get method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		got_blog, err := app.blogpost.Get(id)
		component := base_html.HTML(got_blog.Title, pages.Blogpostviewtempl(partials.BlogPostViewFull(got_blog)), partials.Nav())
		component.Render(context.Background(), f)
	}

}
