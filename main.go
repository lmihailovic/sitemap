package main

import (
	"fmt"
	"github.com/lmihailovic/link"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
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

	baseUrl := strings.Split(url, "/")[0] + "//" + strings.Split(url, "/")[2] + "/"
	domain := strings.Split(url, "/")[2]
	domain = strings.Split(domain, "www.")[1]

	skipExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".svg": true, ".pdf": true, ".doc": true, ".zip": true,
	}

	for path, text := range pageLinks {
		// skips out of domains sites with https protocol
		if strings.HasPrefix(path, "https://") && !strings.Contains(path, domain) {
			//println("skipping: " + path + " is not with " + domain)
			continue
		}
		// skips out of domains sites with http protocol
		if strings.HasPrefix(path, "http://") && !strings.Contains(path, domain) {
			continue
		}

		// skips non-page links
		if strings.HasPrefix(path, "mailto") ||
			strings.HasPrefix(path, "tel") ||
			strings.HasPrefix(path, "#") {
			continue
		}

		// skip non-html files
		ext := strings.ToLower(filepath.Ext(path))
		if skipExt[ext] {
			continue
		}

		// self-explanatory
		if path == ".." {
			continue
		}

		if path == "" {
			println("blank path, link text: " + text)
			continue
		}

		//println("found: " + path)

		if path[0] == '/' {
			path = path[1:]
			path = baseUrl + path
		}

		// skip subdomains
		if !strings.HasPrefix(path, baseUrl) {
			//println("skipping: " + path)
			continue
		}

		if slices.Contains(*visitedPaths, path) {
			continue
		}

		//println("found valid: " + path)
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

	sort.Strings(visited)

	for _, page := range visited {
		fmt.Println(page)
	}

}
