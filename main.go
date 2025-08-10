package main

import (
	"fmt"
	"github.com/lmihailovic/link/parse"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//bodyBytes, err := io.ReadAll(resp.Body)

	body, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	links := link.Parse(body)

	println("Links on page: " + url)
	for k, v := range links {
		fmt.Printf("\npath: %v\ntext: %v\n", k, v)
	}
}
