package scraper

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/attribute"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/test/httpserver"
	"net/http"
	"testing"
	"time"
)

const (
	TestConfigId  = "test"
	TestScraperId = "test_html"
)

func TestNewScraperManager_Start(t *testing.T) {
	url := httpserver.StartTestHttpServer(t)
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to database, error: %s", err)
	}
	m := NewManager(db)
	s := NewScraper(newTestHtmlCrawler(url), getTestHtmlCrawlerCalls(url), nil)
	s.SetTag(attribute.NewTag(TestConfigId, TestScraperId))
	m.RegisterScraper(s)
	m.Start()
	err = db.DeleteMany(database.ScrapedDataTableName, map[string]any{"config_id": "test"})
	if err != nil {
		t.Errorf("Failed to tear down test data, error: %s", err)
	}
}

func newTestHtmlCrawler(url string) *crawler.HtmlCrawler {
	cr := crawler.NewHtmlCrawler(&http.Client{Timeout: 10 * time.Second})
	cr.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	cr.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	return cr
}

func getTestHtmlCrawlerCalls(url string) []*crawler.Call {
	return []*crawler.Call{crawler.NewCall(
		crawler.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/", url), nil),
		crawler.DiscoverRequestType,
	)}
}
