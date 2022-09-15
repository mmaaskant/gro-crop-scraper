package scraper

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/test/httpserver"
	"net/http"
	"testing"
)

func TestScraper_Start(t *testing.T) {
	url := httpserver.StartTestHttpServer(t)
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to database, error: %s", err)
	}
	s := NewScraper(db)
	c := crawler.NewHtmlCrawler("test", &http.Client{})
	c.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	c.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	s.RegisterCrawler(
		c,
		[]*crawler.Call{crawler.NewCrawlerCall(
			fmt.Sprintf("http://%s/", url),
			crawler.DiscoverUrlType,
			http.MethodGet,
			nil,
			nil,
		)},
	)
	s.Start()
	err = db.DeleteMany("scraped_html", map[string]any{"tag": "test"})
	if err != nil {
		t.Errorf("Failed to tear down test data, error: %s", err)
	}
}
