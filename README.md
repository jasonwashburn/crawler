# Crawler

A concurrent web crawler written in Go that crawls websites and generates reports on internal link structure.

## Features

- **Concurrent crawling**: Configurable concurrency level for efficient parallel processing
- **Internal link focus**: Only crawls pages within the same domain as the starting URL
- **URL normalization**: Handles duplicate URLs by normalizing them
- **Link extraction**: Parses HTML to extract all anchor tag links
- **Crawl limits**: Configurable maximum number of pages to crawl
- **Detailed reporting**: Generates sorted reports showing internal link counts

## Usage

```bash
./crawler <URL> <maxConcurrency> <maxPages>
```

### Parameters

- `URL`: The starting URL to begin crawling from
- `maxConcurrency`: Maximum number of concurrent goroutines for crawling (controls parallelism)
- `maxPages`: Maximum number of pages to crawl before stopping

### Example

```bash
./crawler https://example.com 10 100
```

This will crawl `https://example.com` with up to 10 concurrent workers and stop after crawling 100 pages.

## Building

```bash
go build -o crawler
```

## How It Works

1. **Initialization**: Parses command line arguments and sets up the crawler configuration
2. **Concurrent crawling**: Uses goroutines with a semaphore pattern to control concurrency
3. **URL validation**: Only crawls URLs from the same host as the base URL
4. **HTML parsing**: Extracts all `<a href="">` links from each page
5. **Duplicate detection**: Normalizes URLs to avoid crawling the same page multiple times
6. **Progress tracking**: Maintains a thread-safe map of visited pages and their visit counts
7. **Report generation**: Sorts and displays results by link frequency

## Output

The crawler provides real-time logging of its progress and generates a final report showing:

- All discovered internal URLs
- Number of times each URL was referenced by other pages
- Results sorted by reference count (most referenced first)

Example output:

```
=============================
REPORT for https://example.com
=============================
Found 15 internal links to example.com/popular-page
Found 8 internal links to example.com/contact
Found 3 internal links to example.com/about
Found 1 internal links to example.com
```

## Dependencies

- Go 1.23.6+
- `golang.org/x/net/html` for HTML parsing

