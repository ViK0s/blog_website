package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"blog.michalg.net/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	blogpost *models.BlogPostModel
	project  *models.ProjectModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	buildbool := flag.String("bld", "false", "Build or not")

	flag.Parse()

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &app{
		blogpost: &models.BlogPostModel{DB: db},
		project:  &models.ProjectModel{DB: db},
	}

	if *buildbool == "true" {
		app.build()
		log.Println("Built the site, returning")
		return
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /blog/view", app.blogpostview)
	mux.HandleFunc("GET /projects", app.projectshandler)
	mux.HandleFunc("GET /blog", app.blogpostpage)

	log.Println("Starting a server at port :4000")
	log.Fatal(http.ListenAndServe(*addr, mux))
}
