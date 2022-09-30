package main

import (
	"github.com/mmaaskant/gro-crop-scraper/config"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/scraper"
	"log"
)

// main gets all configs and their steps and filters them based on the given flags.
// Config filters are provided as command flags in the format" "--<config_name>",
// Step filters are provided as command flags in the format: "--<step_name>".
// If no filters are provided, all configs and their steps will be executed.
func main() { // TODO: Check CPU and mem usage
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		log.Panicf("Failed to connect to database, error: %s", err)
	}
	sm := scraper.NewManager(db)
	configs := config.GetConfigs()
	for _, c := range configs {
		sm.RegisterScrapers(c.Scrapers)
	}
	sm.Start()
}
