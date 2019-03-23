package main

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

var (
	matchBlogHeader = regexp.MustCompile(`^#\s+(\d+-\d+-\d+):\s+(.*)`)
	subBlogSlug     = regexp.MustCompile(`\W+`)

	posts       = map[time.Time]*BlogPost{}
	sortedPosts = BlogPosts{}

	views   *template.Template
	funcMap = template.FuncMap{
		"date":    func(d time.Time) string { return d.Format("2006-01-02") },
		"xmlTime": func(d time.Time) string { return d.Format("2006-01-03T00:04:05+01:00") },
		"string":  func(t template.HTML) string { return string(t) },
	}
)

type BlogPost struct {
	Slug string

	Language, Title, Body string
	Content               template.HTML

	Created, Updated time.Time

	Published bool
}

func mapBlogPosts() {
	for {
		var err error

		views, err = template.New("views").Funcs(funcMap).ParseGlob("views/*.html")
		if err == nil {
			views, err = views.ParseGlob("views/*.xml")
		}
		if err != nil {
			fmt.Println("error parsing templates:", err)
			fmt.Println("Skipped update:", <-time.After(60*time.Second))
			continue
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
			post.Content = template.HTML(blackfriday.Run([]byte(post.Body)))
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
		post.Content = template.HTML(blackfriday.Run([]byte(post.Body)))
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
