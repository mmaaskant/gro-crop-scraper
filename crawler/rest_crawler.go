package crawler

import (
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"io"
	"log"
	"net/http"
)

// TODO: Write comments

type RestCrawler struct {
	*attributes.Tag
	client *http.Client
}

func NewRestCrawler(id string, c *http.Client) *RestCrawler {
	return &RestCrawler{
		attributes.NewTag("", id),
		c,
	}
}

func (rc *RestCrawler) Crawl(c *Call) *Data {
	b, err := rc.do(c.Request)
	if err != nil {
		log.Printf("Failed to crawl url: %s, error: %s", c.URL.String(), err)
		return NewData(rc.Tag, c, "", nil, err)
	}
	return NewData(rc.Tag, c, b, nil, err)
}

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
