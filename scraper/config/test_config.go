package config

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"net/http"
)

const (
	ScraperConfigTestTag        = "test"
	ScraperConfigTestHtmlOrigin = "test_html"
)

func NewTestConfig(testServerUrl string) *Config {
	c := newConfig(ScraperConfigTestTag)
	c.AddCrawler(newTestHtmlCrawler(testServerUrl))
	return c
}

func newTestHtmlCrawler(url string) (*crawler.HtmlCrawler, []*crawler.Call) {
	cr := crawler.NewHtmlCrawler(ScraperConfigTestHtmlOrigin, &http.Client{})
	cr.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	cr.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	return cr, []*crawler.Call{crawler.NewCrawlerCall(
		fmt.Sprintf("http://%s/", url),
		crawler.DiscoverUrlType,
		http.MethodGet,
		nil,
		nil,
	)}
}
