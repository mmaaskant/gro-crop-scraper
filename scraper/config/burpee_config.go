package config

import (
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"net/http"
	"time"
)

const (
	ScraperConfigBurpeeOrigin      = "burpee"
	ScraperConfigBurpeeHtmlDataId  = "burpee_html"
	ScraperConfigBurpeeZonesDataId = "burpee_zones"
)

// NewBurpeeConfig returns a new instance of *Config holding all components required to scrape the Burpee supplier sources.
func NewBurpeeConfig() *Config {
	c := newConfig(ScraperConfigBurpeeOrigin)
	c.AddCrawler(newBurpeeHtmlCrawler())
	return c
}

// newBurpeeHtmlCrawler returns an instance of crawler.HtmlCrawler configured to scrape
// the Burpee website and a slice of crawler.Call instances to kick off the crawling process.
func newBurpeeHtmlCrawler() (*crawler.HtmlCrawler, []*crawler.Call) {
	cr := crawler.NewHtmlCrawler(ScraperConfigBurpeeHtmlDataId, &http.Client{Timeout: 90 * time.Second})
	cr.AddDiscoveryUrlRegex(`(https?:\/\/)?www\.burpee\.com\/?(vegetables|flowers|perennials|herbs|fruit)([\w\/-]*)(\?p=\d{1,3})?(&is_scroll=1)?`)
	cr.AddExtractUrlRegex(`(https?:\/\/)?www\.burpee\.com\/([\w\-]*)(prod\d*.html)(\/)?`)
	return cr, []*crawler.Call{crawler.NewCall(
		crawler.NewRequest(http.MethodGet, "https://www.burpee.com", nil),
		crawler.DiscoverRequestType,
	)}
}
