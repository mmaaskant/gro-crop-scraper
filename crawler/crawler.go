package crawler

import (
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"io"
	"log"
	"net/http"
)

const (
	DiscoverRequestType string = "DISCOVER"
	ExtractRequestType  string = "EXTRACT"
)

// Crawler crawls any URL and returns an instance of Data containing what it has found.
type Crawler interface {
	attributes.Taggable
	Crawl(c *Call) *Data
}

// Call contains everything that is required by Crawler to make a request.
type Call struct {
	*http.Request
	RequestType string
}

// NewCrawlerCall returns a new instance of Call.
func NewCrawlerCall(r *http.Request, RequestType string) *Call {
	return &Call{
		r,
		RequestType,
	}
}

// NewRequest returns a new instance of http.Request and logs a fatal error if it fails.
func NewRequest(method string, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("Failed to create crawler HTTP request, error: %s", err)
	}
	return r
}

// Data contains all data that was found by a Crawler.Crawl call, the Call itself, and a collection of found calls.
type Data struct {
	*attributes.Tag
	Call       *Call
	Data       string
	FoundCalls []*Call
	Error      error
}

// NewCrawlerData returns a new instance of Data
func NewCrawlerData(t *attributes.Tag, call *Call, data string, foundCalls []*Call, err error) *Data {
	return &Data{
		t,
		call,
		data,
		foundCalls,
		err,
	}
}
