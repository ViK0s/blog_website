package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"

	base_html "blog.michalg.net/ui/html"
	"blog.michalg.net/ui/html/pages"
	"blog.michalg.net/ui/html/partials"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// changelinks handles the link changing and rerendering of the html file
// this isnt the best way of implementing this probably, but it's mostly reusable code so I made it into a function to use
// in every builder function
func changelinks(ff *os.File, links map[string]string, filename string) {
	//get the html elements from the file
	//some name changes need to be made here because calling everything doc is stupid
	doc, err := html.Parse(ff)
	if err != nil {
		log.Fatalf("failed to parse html %v", err)
	}

	doc_query := goquery.NewDocumentFromNode(doc)
	//find every <a> tag and check the href attr, change if found on the change list
	doc_query.Find("a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists {

			//okay so because for the blogs we use /blog/view?=id we know the lenght of the link will be 15
			//we need to check that so that we know wecan slice the string
			//we slice the string to know that it is the one we need
			//and then we modify it
			if len(href) == 15 {
				if href[0:10] == "/blog/view" {
					item.SetAttr("href", "./viewid="+href[len(href)-1:]+".html")
				}
			}
			// change the link by using a slice
			// for now this doesnt do anything if the link found isn't on the list
			// but I'm pretty sure there should be something done about this
			_, ok := links[href]
			if ok {
				item.SetAttr("href", links[href])
			}
		}
	})

	//this part is needed because sometimes we want to change the css links inside the <head> tag
	doc_query.Find("link").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists {

			// change the link by using a slice
			// for now this doesnt do anything if the link found isn't on the list
			// but I'm pretty sure there should be something done about this
			_, ok := links[href]
			if ok {
				item.SetAttr("href", links[href])
			}
		}
	})
	//part for changing img links
	doc_query.Find("img").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("src")
		if exists {

			item.SetAttr("src", "../ui/static/img/"+href[12:])

		}
	})

	if err := os.Truncate(filename, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	doc_query.Each(func(index int, item *goquery.Selection) {
		err = goquery.Render(ff, item)
	})

}

// passing the pointer to app but not expanding the app struct because I want the app to not have random shit in it
func homebuild(app *app, links map[string]string) {
	f, err := os.Create("buildstat/home.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	component := base_html.HTML("Home", pages.Hometmp(partials.LatestBlogPostTempl(app.blogpost.Latest()[0]), partials.Another_Section(app.blogpost.Latest()[1])), partials.Nav())
	component.Render(context.Background(), f)

	f.Close()

	ff, err := os.OpenFile("buildstat/home.html", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed to open specified file %v", err)
	}

	changelinks(ff, links, "buildstat/home.html")
}

func buildabout(links map[string]string) {
	f, err := os.Create("buildstat/about.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	component := base_html.HTML("About", pages.Abouttmp(partials.Contact_Section()), partials.Nav())
	component.Render(context.Background(), f)

	f.Close()

	ff, err := os.OpenFile("buildstat/about.html", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed to open specified file %v", err)
	}

	changelinks(ff, links, "buildstat/about.html")
}

func blogviewbuild(app *app, links map[string]string) {
	kros := app.blogpost.Latest()[0].ID
	var f *os.File
	var err error
	// start building blogposts
	for id := range kros + 1 {
		if id > 0 {
			f, err = os.Create("buildstat/viewid=" + strconv.Itoa(id) + ".html")

			if err != nil {
				log.Fatalf("Failed at creating a file %v", err)
				return
			}

			got_blog, err := app.blogpost.Get(id)
			if err != nil {
				log.Fatalf("Couldn't get a blogpost from the database %v", err)
				return
			}
			component := base_html.HTML(got_blog.Title, pages.Blogpostviewtempl(partials.BlogPostViewFull(got_blog)), partials.Nav())
			component.Render(context.Background(), f)

			f.Close()

			ff, err := os.OpenFile("buildstat/viewid="+strconv.Itoa(id)+".html", os.O_RDWR, os.ModeAppend)
			if err != nil {
				log.Fatalf("failed to open specified file %v", err)
			}

			changelinks(ff, links, "buildstat/viewid="+strconv.Itoa(id)+".html")
		}

	}
}

func buildprojects(app *app, links map[string]string) {

	// first create the html file
	f, err := os.Create("buildstat/projects.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	//render to the file
	component := base_html.HTML("Projects", pages.ProjectsSite(partials.Projecttemp(app.project.Latest())), partials.Nav())
	component.Render(context.Background(), f)

	f.Close()

	// because the f pointer doesnt really update we need to reopen the file
	// we use the openfile function because we need to specify that the file will be used to read and write
	// this can also be put inside the changelinks function tbh
	ff, err := os.OpenFile("buildstat/projects.html", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed to open specified file %v", err)
	}

	changelinks(ff, links, "buildstat/projects.html")

}

func buildblogs(app *app, links map[string]string) {
	f, err := os.Create("buildstat/blogs.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	component := base_html.HTML("Blogs", pages.Blogsite(partials.BlogPostFull(app.blogpost.Latest()), partials.Archive(app.blogpost.Latest())), partials.Nav())
	component.Render(context.Background(), f)

	ff, err := os.OpenFile("buildstat/blogs.html", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed to open specified file %v", err)
	}

	changelinks(ff, links, "buildstat/blogs.html")
}

// every handler should probably have it's own function so that it's easier to read
// and that every handler has it's own way to handle errors
// but for now I left it as is beacuse I'm too lazy

func (app *app) build() {

	newpath := filepath.Join(".", "buildstat")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	if err != nil {
		log.Print(err)
		return
	}

	//define a map so that we know which links to change
	links := map[string]string{"/static/css/main.css": "./main.css", "/": "./home.html", "/projects": "./projects.html", "/blog": "./blogs.html", "/about": "./about.html"}

	data, err := os.ReadFile("./ui/static/css/main.css")
	if err != nil {
		log.Print(err)
		return
	}
	// Write data to dst
	err = os.WriteFile(newpath+"/main.css", data, 0644)
	if err != nil {
		log.Print(err)
		return
	}

	fff, err := os.OpenFile("buildstat/main.css", os.O_RDWR, os.ModeAppend)
	scannery := bufio.NewScanner(fff)

	var data2 string
	for scannery.Scan() {
		if len(scannery.Text()) > 12 {
			if scannery.Text()[0:11] == "        url" {
				data2 += "        url(" + "../ui/" + scannery.Text()[13:]
				continue
			}
		}
		data2 += scannery.Text()
	}

	if err := os.Truncate("buildstat/main.css", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	os.WriteFile("buildstat/main.css", []byte(data2), 0644)

	homebuild(app, links)
	buildabout(links)
	blogviewbuild(app, links)
	buildprojects(app, links)
	buildblogs(app, links)

}
