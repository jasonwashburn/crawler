package main

import (
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
					if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
						links = append(links, link)
					} else if strings.HasPrefix(link, "/") {
						baseURLParts := strings.Split(rawBaseURL, "/")
						baseURL := strings.Join(baseURLParts[:3], "/")
						links = append(links, baseURL+link)
					} else if strings.HasPrefix(link, "../") {
						baseURLParts := strings.Split(rawBaseURL, "/")
						baseURL := strings.Join(baseURLParts[:3], "/")
						link = strings.TrimPrefix(link, "../")
						link = baseURL + "/" + link
						links = append(links, link)
					}
				}
			}
		}
	}
	return links, nil
}
