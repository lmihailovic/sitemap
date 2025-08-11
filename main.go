package main

import (
	"fmt"
	"github.com/lmihailovic/link/parse"
	"golang.org/x/net/html"
	"net/http"
	"os"
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

func GetSitePages(url string) ([]string, error) {
	body, err := GetPageHtml(url)
	if err != nil {
		return nil, err
	}

	pageLinks := link.Parse(body)
	paths := make([]string, 0)

	for path, _ := range pageLinks {
		if strings.HasPrefix(path, "mailto") ||
			strings.HasPrefix(path, "tel") ||
			strings.HasPrefix(path, "#") {
			continue
		} else if strings.Contains(path, ".") {
			continue
		}
		subUrl := url + path
		paths = append(paths, subUrl)
	}

	return paths, nil
}

func main() {
	url := os.Args[1]

	pages, err := GetSitePages(url)
	if err != nil {
		panic(err)
	}

	for _, page := range pages {
		fmt.Println(page)
	}

}
