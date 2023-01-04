package md

import "fmt"

const (
	// ExampleURL indicates the default value used when an URL is not explicitly defined
	ExampleURL = "http://example.com"
)

// MakeLink produces a link in proper Markdown syntax
func MakeLink(url string, caption string) string {
	if url == "" {
		url = ExampleURL
	}
	if caption == "" {
		return url
	}

	return fmt.Sprintf("[%s](%s)", caption, url)
}
