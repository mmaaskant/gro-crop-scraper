package crawler

import (
	"github.com/mmaaskant/gro-crop-scraper/attributes"
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
	return nil
}
