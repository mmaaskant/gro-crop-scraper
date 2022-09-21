package crawler

import (
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"io"
	"log"
	"net/http"
)

// RestCrawler crawls REST APIs using the provided Call instance.
type RestCrawler struct {
	*attributes.Tag
	client *http.Client
}

// NewRestCrawler returns a new instance of RestCrawler.
func NewRestCrawler(c *http.Client) *RestCrawler {
	return &RestCrawler{
		nil,
		c,
	}
}

func (rc *RestCrawler) SetTag(t *attributes.Tag) {
	rc.Tag = t
}

// Crawl starts crawling based on the given Call instance and returns a Data instance
// containing the response as a string and any other relevant data found along the way.
func (rc *RestCrawler) Crawl(c *Call) *Data {
	b, err := rc.do(c.Request)
	if err != nil {
		log.Printf("Failed to crawl url: %s, error: %s", c.URL.String(), err)
		return NewData(rc.Tag, c, "", nil, err)
	}
	return NewData(rc.Tag, c, b, nil, err)
}

// do calls the provided http.Request, stringifies the result body and returns it.
func (rc *RestCrawler) do(req *http.Request) (string, error) {
	resp, err := rc.client.Do(req)
	if err != nil {
		log.Printf("Failed to call url: %s, error: %s", req.URL.String(), err)
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), err
}
