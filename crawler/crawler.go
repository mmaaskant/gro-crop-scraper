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

// Crawler crawls any URL and returns Data containing what it has found,
// it also implements attributes.Taggable allowing it to tag said Data.
type Crawler interface { // TODO: Clients should support retries
	attributes.Taggable
	Crawl(c *Call) *Data
}

// Call wraps around a http.Request and adds a RequestType which should be either DiscoverRequestType or ExtractRequestType.
// In which the former will be used only to discover new URLs, and the latter will be stored locally for further processing.
type Call struct {
	*http.Request
	RequestType string
}

func NewCall(r *http.Request, RequestType string) *Call {
	return &Call{
		r,
		RequestType,
	}
}

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

func NewData(t *attributes.Tag, call *Call, data string, foundCalls []*Call, err error) *Data {
	return &Data{
		t,
		call,
		data,
		foundCalls,
		err,
	}
}
