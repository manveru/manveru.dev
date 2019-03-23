package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux" // the 800 pound gorilla
)

var (
	flagHost   = flag.String("host", "", "the host to listen on")
	flagPort   = flag.String("port", "8000", "the port to listen on")
	flagPublic = flag.String("public", "./public", "serve static files from this directory")
)

func main() {
	// parse command line options
	flag.Parse()

	go mapBlogPosts()

	setupServer()
}

func setupServer() {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/", mainIndex)
	r.HandleFunc("/contact", mainContact)
	r.HandleFunc("/archive", blogArchive)
	r.HandleFunc(`/blog/show/{id:\d{4}-\d{2}-\d{2}}/{language}/{slug}`, blogShow)
	r.HandleFunc(`/blog/feed.atom`, blogAtom)

	// static files
	for _, name := range []string{"css", "js", "img", "static", "google-code-prettify"} {
		prefix, dirName := "/"+name+"/", *flagPublic+"/"+name
		r.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(http.Dir(dirName))))
	}

	err := http.ListenAndServe(*flagHost+":"+*flagPort, r)

	if err != nil {
		fmt.Println("Could not start server:", err)
	}
}

func render(out http.ResponseWriter, name string, val interface{}) {
	err := views.ExecuteTemplate(out, name, val)

	if err != nil {
		panic(err)
	}
}

type args map[string]interface{}

func notFound(out http.ResponseWriter, req *http.Request) {
	render(out, "404", args{
		"Title":   "Oops",
		"Subject": "[manveru.dev] 404 @ " + req.RequestURI,
		"Body":    fmt.Sprintf("%#v", req),
	})
}

func mainIndex(out http.ResponseWriter, req *http.Request) {
	render(out, "index", args{"Title": "Home", "Posts": sortedPostsN(20)})
}

func mainContact(out http.ResponseWriter, req *http.Request) {
	render(out, "contact", args{"Title": "Contact"})
}

func blogArchive(out http.ResponseWriter, req *http.Request) {
	render(out, "blog/archive", args{"Title": "Archive", "Posts": sortedPostsN(1000)})
}

func blogShow(out http.ResponseWriter, req *http.Request) {
	post := postById(mux.Vars(req)["id"])
	render(out, "blog/show", args{"Title": post.Title, "Post": post})
}

func blogAtom(out http.ResponseWriter, req *http.Request) {
	out.Header().Set("Content-Type", "application/atom+xml")
	posts := sortedPostsN(100)
	render(out, "blog/atom", args{
		"Updated": posts[0].Updated,
		"Posts":   posts,
	})
}
