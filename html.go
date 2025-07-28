package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var links []string
	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					if !strings.HasPrefix(link, "http") {
						// TODO: Handle relative URLs with ..
						// TODO: Handle relative URLs starting with root
						link = fmt.Sprintf("%s/%s", rawBaseURL, strings.TrimLeft(link, "/"))
					}
					links = append(links, link)
				}
			}
		}
	}
	return links, nil
}
