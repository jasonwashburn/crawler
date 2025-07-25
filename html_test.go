package main

import (
	"slices"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "nested links and mixed formats",
			inputURL: "https://example.com",
			inputBody: `
<html>
    <head>
        <title>Test Page</title>
    </head>
    <body>
        <div>
            <p>Check out <a href="/docs">our documentation</a></p>
            <ul>
                <li><a href="https://github.com/user/repo">GitHub</a></li>
                <li><a href="../parent/page.html">Parent page</a></li>
            </ul>
        </div>
        <footer>
            <a href="/about">About Us</a>
        </footer>
    </body>
</html>
`,
			expected: []string{
				"https://example.com/docs",
				"https://github.com/user/repo",
				"https://example.com/parent/page.html",
				"https://example.com/about",
			},
		},
		{
			name:     "multiple links same domain",
			inputURL: "https://example.com/blog",
			inputBody: `
<html>
    <body>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
            <a href="/contact">Contact</a>
        </nav>
        <main>
            <a href="/blog/post-1">First Post</a>
            <a href="/blog/post-2">Second Post</a>
        </main>
    </body>
</html>
`,
			expected: []string{
				"https://example.com/",
				"https://example.com/about",
				"https://example.com/contact",
				"https://example.com/blog/post-1",
				"https://example.com/blog/post-2",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !slices.Equal(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URLs: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
