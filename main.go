package main

import (
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/scraper"
	"github.com/mmaaskant/gro-crop-scraper/scraper/config"
	"log"
)

// main gets all registered configs and filters them accordingly based on the given flags,
// filters are provided as command flags in the format "--<origin_name>".
// If no filters are provided, all configs are ran through an instance of scraper.Scraper.
func main() {
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		log.Fatalf("Failed to connect to database, error: %s", err)
	}
	s := scraper.NewScraper(db)
	configs := config.GetRegisteredConfigs()
	for _, c := range configs {
		s.RegisterConfig(c)
	}
	s.Start()
}
