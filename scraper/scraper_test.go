package scraper

import (
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/scraper/config"
	"github.com/mmaaskant/gro-crop-scraper/test/httpserver"
	"testing"
)

func TestScraper_Start(t *testing.T) {
	url := httpserver.StartTestHttpServer(t)
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to database, error: %s", err)
	}
	s := NewScraper(db)
	s.RegisterConfig(config.NewTestConfig(url))
	s.Start()
	err = db.DeleteMany(database.DbScrapedDataTableName, map[string]any{"origin": "test"})
	if err != nil {
		t.Errorf("Failed to tear down test data, error: %s", err)
	}
}
