package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args

	if len(args) < 4 {
		fmt.Println("Usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	if len(args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL, err := url.Parse(args[1])
	if err != nil {
		fmt.Println("error parsing base URL:", err)
		os.Exit(1)
	}

	fmt.Println("starting crawl of: ", baseURL.String())

	pages := make(map[string]int)
	maxConcurrency, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("error parsing maxConcurrency:", err)
	}
	maxPages, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("error parsing maxPages:", err)
	}

	cfg := config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	cfg.crawlPage(baseURL.String())

	cfg.wg.Wait()
	for page, count := range pages {
		fmt.Printf("Page: %s, Count: %d\n", page, count)
	}
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		fmt.Println("maximum number of pages reached:", cfg.maxPages)
		return
	}
	cfg.mu.Unlock()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("error parsing current URL:", err)
		return
	}

	if currentURL.Host != cfg.baseURL.Host {
		fmt.Println("skipping external link:", currentURL.String())
		return
	}

	normalizedURL, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println("error normalizing URL:", err)
		return
	}

	if isFirst := cfg.addPageVisit(normalizedURL); !isFirst {
		fmt.Println("already crawled:", normalizedURL, "count:", cfg.pages[normalizedURL])
	} else {
		fmt.Println("crawling:", currentURL.String())
		html, err := getHTML(currentURL.String())
		if err != nil {
			fmt.Println("error fetching HTML for", currentURL.String(), ":", err)
			return
		}

		urls, err := getURLsFromHTML(html, cfg.baseURL.String())
		if err != nil {
			fmt.Println("error extracting links from", normalizedURL, ":", err)
			return
		}

		fmt.Printf("Extracted %d links from %s\n", len(urls), normalizedURL)

		for _, link := range urls {
			cfg.wg.Add(1)
			go func(link string) {
				fmt.Println("spawning goroutine for link:", link)
				defer cfg.wg.Done()
				cfg.concurrencyControl <- struct{}{} // Acquire a slot in the channel
				cfg.crawlPage(link)
				<-cfg.concurrencyControl // Release the slot in the channel
			}(link)
		}
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, exists := cfg.pages[normalizedURL]
	isFirst = !exists
	fmt.Println("adding page visit for:", normalizedURL, "isFirst:", isFirst)
	cfg.pages[normalizedURL]++
	return
}
