package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	normalizedURL := strings.TrimRight(fmt.Sprint(url.Host, url.Path), "/")
	return normalizedURL, nil
}
