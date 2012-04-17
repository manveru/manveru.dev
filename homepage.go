package main

import (
	"bufio"
	"code.google.com/p/gorilla/mux" // the 800 pound gorilla
	"errors"
	"flag"
	"fmt"
	"github.com/russross/blackfriday" // markdown
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	flagHost   = flag.String("host", "", "the host to listen on")
	flagPort   = flag.String("port", "8000", "the port to listen on")
	flagPublic = flag.String("public", "./public", "serve static files from this directory")

	matchBlogHeader = regexp.MustCompile(`^#\s+(\d+-\d+-\d+):\s+(.*)`)
	subBlogSlug     = regexp.MustCompile(`\W+`)

	views       *template.Template
	posts       = map[time.Time]*BlogPost{}
	sortedPosts = BlogPosts{}
	funcMap     = template.FuncMap{"date": formatDate}
)

func main() {
	// parse command line options
	flag.Parse()

	go mapBlogPosts()

	setupServer()
}

func formatDate(d time.Time) string {
	return d.Format("2006-01-02")
}

func mapBlogPosts() {
	for {
		var err error
		views, err = template.New("views").Funcs(funcMap).ParseGlob("views/*.html")
		if err != nil {
			panic(err)
		}

		filepath.Walk("posts", func(path string, info os.FileInfo, err error) error {
			if err == nil {
				if filepath.Ext(path) == ".md" {
					return loadBlogPosts(path, info)
				}
			}
			return nil
		})

		sortedPosts = BlogPosts{}

		for _, post := range posts {
			sortedPosts = append(sortedPosts, post)
		}

		sort.Sort(sortedPosts)

		fmt.Println("Updating homepage:", <-time.After(60*time.Second))
	}
}

type BlogPosts []*BlogPost

func (bp BlogPosts) Len() int           { return len(bp) }
func (bp BlogPosts) Swap(a, b int)      { bp[a], bp[b] = bp[b], bp[a] }
func (bp BlogPosts) Less(a, b int) bool { return bp[a].Created.After(bp[b].Created) }

func sortedPostsN(limit int) BlogPosts {
	if limit > len(sortedPosts) {
		limit = len(sortedPosts)
	}

	return sortedPosts[:limit]
}

func postById(id string) *BlogPost {
	postCreated, err := time.Parse("2006-01-02", id)
	if err != nil {
		panic(err)
	}

	return posts[postCreated]
}

type BlogPost struct {
	Slug string

	Language, Title, Body string
	Content               template.HTML

	Created, Updated time.Time

	Published bool
}

func loadBlogPosts(path string, info os.FileInfo) (err error) {
	tmpPosts := []*BlogPost{}
	lines := []string{}
	var post *BlogPost

	eachLine(path, func(line string, err error) {
		matches := matchBlogHeader.FindStringSubmatch(line)
		if len(matches) != 3 {
			lines = append(lines, line)
			return
		}

		if post != nil {
			post.Body = strings.Join(lines, "\n")
			post.Content = template.HTML(blackfriday.MarkdownCommon([]byte(post.Body)))
			tmpPosts = append(tmpPosts, post)
			lines = lines[:0]
		}

		created, err := time.Parse("2006-01-02", matches[1])
		if err != nil {
			return
		}

		post = &BlogPost{
			Created:   created,
			Updated:   info.ModTime(),
			Published: true,
			Language:  "en",
			Title:     matches[2],
			Slug:      strings.Trim(subBlogSlug.ReplaceAllString(matches[2], "-"), "-"),
			Body:      "",
		}
	})

	if post != nil {
		post.Body = strings.Join(lines, "\n")
		post.Content = template.HTML(blackfriday.MarkdownCommon([]byte(post.Body)))
		tmpPosts = append(tmpPosts, post)
	}

	for _, post := range tmpPosts {
		posts[post.Created] = post
	}

	return
}

func eachLine(path string, f func(line string, err error)) {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Println("os.Open", path, err)
		return
	}

	// lines longer than 10k bytes will be punished
	fdReader := bufio.NewReaderSize(fd, 10000)

	var line []byte
	isPrefix := false
	for {
		line, isPrefix, err = fdReader.ReadLine()

		if err != nil {
			if err != io.EOF {
				f("", err)
			}
			break
		}

		if isPrefix {
			err = errors.New("line too long")
		}

		f(string(line), err)
	}
}

func setupServer() {
	http.Handle("/", http.FileServer(http.Dir("/tmp")))

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/", mainIndex)
	r.HandleFunc("/contact", mainContact)
	r.HandleFunc("/archive", blogArchive)
	r.HandleFunc(`/blog/show/{id:\d{4}-\d{2}-\d{2}}/{language}/{slug}`, blogShow)
	for _, name := range []string{"css", "js", "img", "static"} {
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

type rargs map[string]interface{}
type indexArgs struct {
	Title string
	Posts BlogPosts
}

func notFound(out http.ResponseWriter, req *http.Request) {
	render(out, "404", rargs{"Title": "Oops"})
}

func mainIndex(out http.ResponseWriter, req *http.Request) {
	render(out, "index", indexArgs{Title: "Home", Posts: sortedPostsN(20)})
}

func mainContact(out http.ResponseWriter, req *http.Request) {
	render(out, "contact", rargs{"Title": "Contact"})
}

func blogArchive(out http.ResponseWriter, req *http.Request) {
	render(out, "blog/archive", rargs{"Title": "Archive", "Posts": sortedPostsN(100)})
}

func blogShow(out http.ResponseWriter, req *http.Request) {
	post := postById(mux.Vars(req)["id"])
	render(out, "blog/show", rargs{"Title": post.Title, "Post": post})
}
