package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	var links []string
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var visitNode func(*html.Node)
	visitNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					resolvedURL := baseURL.ResolveReference(link)
					links = append(links, resolvedURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visitNode(c)
		}
	}

	visitNode(doc)
	return links, nil
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to fetch URL %s: %s", rawURL, resp.Status)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("URL %s is not an HTML document", rawURL)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
