package config

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"net/http"
)

const (
	ScraperConfigTestOrigin     = "test"
	ScraperConfigTestHtmlDataId = "test_html"
)

func NewTestConfig(testServerUrl string) *Config {
	c := newConfig(ScraperConfigTestOrigin)
	c.AddCrawler(newTestHtmlCrawler(testServerUrl))
	return c
}

func newTestHtmlCrawler(url string) (*crawler.HtmlCrawler, []*crawler.Call) {
	cr := crawler.NewHtmlCrawler(ScraperConfigTestHtmlDataId, &http.Client{})
	cr.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	cr.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	return cr, []*crawler.Call{crawler.NewCrawlerCall(
		crawler.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/", url), nil),
		crawler.DiscoverRequestType,
	)}
}
