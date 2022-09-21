package config

import (
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/scraper"
	"net/http"
	"net/url"
	"time"
)

const (
	BurpeeConfigId      = "burpee"
	BurpeeHtmlScraperId = "burpee_html"
	BurpeeZoneScraperId = "burpee_zones"
)

// NewBurpeeConfig holds components configured to scrape https://burpee.com.
func NewBurpeeConfig() *Config {
	c := newConfig(BurpeeConfigId)
	c.AddScraper(BurpeeHtmlScraperId, scraper.NewScraper(newBurpeeHtmlCrawler(), getBurpeeHtmlCrawlerCalls(), nil))
	c.AddScraper(BurpeeZoneScraperId, scraper.NewScraper(newBurpeeZonesCrawler(), getBurpeeZonesCrawlerCalls(), nil))
	return c
}

// newBurpeeHtmlCrawler returns an instance of crawler.HtmlCrawler configured to scrape
// the Burpee website and a slice of crawler.Call instances to kick off the crawling process.
func newBurpeeHtmlCrawler() *crawler.HtmlCrawler {
	cr := crawler.NewHtmlCrawler(&http.Client{Timeout: 90 * time.Second})
	cr.AddDiscoveryUrlRegex(`(https?:\/\/)?www\.burpee\.com\/?(vegetables|flowers|perennials|herbs|fruit)([\w\/-]*)(\?p=\d{1,3})?(&is_scroll=1)?`)
	cr.AddExtractUrlRegex(`(https?:\/\/)?www\.burpee\.com\/([\w\-]*)(prod\d*.html)(\/)?`)
	return cr
}

func getBurpeeHtmlCrawlerCalls() []*crawler.Call {
	return []*crawler.Call{crawler.NewCall(
		crawler.NewRequest(http.MethodGet, "https://www.burpee.com", nil),
		crawler.DiscoverRequestType,
	)}
}

// newBurpeeZonesCrawler returns an instance of crawler.HtmlCrawler configured to scrape
// the Burpee growing zone and crop category attribute API.
func newBurpeeZonesCrawler() *crawler.RestCrawler {
	cr := crawler.NewRestCrawler(&http.Client{Timeout: 30 * time.Second})
	return cr
}

func getBurpeeZonesCrawlerCalls() []*crawler.Call {
	calls := make([]*crawler.Call, 0)
	for _, r := range getGrowingZoneRequests() {
		c := crawler.NewCall(r, crawler.ExtractRequestType)
		calls = append(calls, c)
	}
	return calls
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

// TODO: Should these be pulled from a file instead?
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
