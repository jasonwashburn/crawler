package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[1]
	fmt.Println("starting crawl of: ", baseURL)

	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)

	for page, count := range pages {
		fmt.Printf("Page: %s, Count: %d\n", page, count)
	}
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println("error parsing base URL:", err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("error parsing current URL:", err)
		return
	}

	if currentURL.Host != baseURL.Host {
		fmt.Println("skipping external link:", currentURL.String())
		return
	}

	normalizedURL, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println("error normalizing URL:", err)
		return
	}

	if _, exists := pages[normalizedURL]; exists {
		pages[normalizedURL]++
		fmt.Println("already crawled:", normalizedURL, "count:", pages[normalizedURL])
	} else {
		pages[normalizedURL] = 1
		fmt.Println("crawling:", currentURL.String())
		html, err := getHTML(currentURL.String())
		if err != nil {
			fmt.Println("error fetching HTML for", currentURL.String(), ":", err)
			return
		}

		urls, err := getURLsFromHTML(html, rawBaseURL)
		if err != nil {
			fmt.Println("error extracting links from", normalizedURL, ":", err)
			return
		}

		fmt.Printf("Extracted %d links from %s\n", len(urls), normalizedURL)

		for _, link := range urls {
			crawlPage(rawBaseURL, link, pages)
		}
	}
}
