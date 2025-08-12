package main

import (
	"fmt"
	"github.com/lmihailovic/link/parse"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"slices"
	"strings"
)

/*
	SAMPLE SITEMAP XML

	<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
		<url>
			<loc>http://www.example.com/</loc>
		</url>
		<url>
			<loc>http://www.example.com/dogs</loc>
		</url>
	</urlset>
*/

func GetPageHtml(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetSitePages(url string, visitedPaths *[]string) error {
	body, err := GetPageHtml(url)
	if err != nil {
		return err
	}

	pageLinks := link.Parse(body)
	foundPaths := make([]string, 0)

	fmt.Printf("\nHit page %v\n", url)

	for path, _ := range pageLinks {
		if strings.HasPrefix(path, "https://") ||
			strings.HasPrefix(path, "http://") ||
			strings.HasPrefix(path, "mailto") ||
			strings.HasPrefix(path, "tel") ||
			strings.HasPrefix(path, "#") ||
			(strings.Contains(path, ".") && !strings.HasSuffix(path, ".html")) {
			continue
		}

		baseUrl := strings.Split(url, "/")[0] + "//" + strings.Split(url, "/")[2] + "/"

		if path[0] == '/' {
			path = path[1:]
			path = baseUrl + path
		}

		if slices.Contains(*visitedPaths, path) {
			continue
		}

		if path == ".." {
			continue
		}

		//println("found: " + path)
		foundPaths = append(foundPaths, path)
	}

	*visitedPaths = append(*visitedPaths, url)
	for _, path := range foundPaths {
		if slices.Contains(*visitedPaths, path) {
			continue
		}

		err := GetSitePages(path, visitedPaths)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	url := os.Args[1]

	var visited []string
	err := GetSitePages(url, &visited)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found pages:")

	for _, page := range visited {
		fmt.Println(page)
	}

}
