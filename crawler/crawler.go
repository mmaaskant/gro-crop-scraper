package crawler

import "github.com/mmaaskant/gro-crop-scraper/attributes"

const (
	DiscoverUrlType string = "DISCOVER"
	ExtractUrlType  string = "EXTRACT"
)

// Crawler crawls any URL and returns an instance of Data containing what it has found.
type Crawler interface {
	attributes.Taggable
	Crawl(c *Call) *Data
}

// Call contains everything that is required by Crawler to make a request.
type Call struct {
	Url     string
	UrlType string
	Method  string
	Headers map[string]string
	Params  map[string]string
}

// NewCrawlerCall returns a new instance of Call.
func NewCrawlerCall(url string, UrlType string, method string, headers map[string]string, params map[string]string) *Call {
	return &Call{
		url,
		UrlType,
		method,
		headers,
		params,
	}
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
