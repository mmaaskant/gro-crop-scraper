package config

import (
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"net/http"
	"net/url"
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

// newBurpeeZonesCrawler returns an instance of crawler.HtmlCrawler configured to scrape
// the Burpee growing zone and crop category attribute API.
func newBurpeeZonesCrawler() (*crawler.RestCrawler, []*crawler.Call) {
	calls := make([]*crawler.Call, 0)
	for _, r := range getGrowingZoneRequests() {
		c := crawler.NewCall(r, crawler.ExtractRequestType)
		calls = append(calls, c)
	}
	cr := crawler.NewRestCrawler(ScraperConfigBurpeeZonesDataId, &http.Client{Timeout: 30 * time.Second})
	return cr, calls
}

// getGrowingZoneRequests compiles a slice of http.Request for each growing zone that is available
// within the Burpee API.
func getGrowingZoneRequests() []*http.Request {
	requests := make([]*http.Request, 0)
	for _, zip := range getGrowingZoneZipcodes() {
		r := crawler.NewRequest(http.MethodGet, "https://www.burpee.com/location/index/index", nil)
		d := url.Values{"zipcode": []string{zip}}
		r.URL.RawQuery = d.Encode()
		requests = append(requests, r)
	}
	return requests
}

// getGrowingZoneZipcodes returns a map of growing zones and a zipcode within said zone.
func getGrowingZoneZipcodes() map[int]string {
	return map[int]string{
		1:  "99722",
		2:  "99731",
		3:  "99736",
		4:  "59317",
		5:  "57785",
		6:  "97867",
		7:  "98815",
		8:  "98589",
		9:  "70517",
		10: "34104",
		11: "33037",
		12: "96778",
		13: "96863",
	}
}
