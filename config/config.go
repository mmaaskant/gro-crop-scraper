package config

import (
	"flag"
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"github.com/mmaaskant/gro-crop-scraper/scraper"
)

var configs = []*Config{
	NewBurpeeConfig(),
}

func init() {
	handleFlags()
}

// Config groups scraper.Scraper configurations under an ID,
// this ID is passed along scraper.Scraper, its components and any data that it handles.
type Config struct {
	Id       string
	Scrapers []*scraper.Scraper
}

func newConfig(id string) *Config {
	return &Config{
		id,
		make([]*scraper.Scraper, 0),
	}
}

// GetConfigs returns all Config instances that have been flagged, or all of them if none have been flagged.
func GetConfigs() []*Config {
	rc := make([]*Config, 0)
	for _, c := range configs {
		if flagToBool(flag.Lookup(c.Id)) == true {
			rc = append(rc, c)
		}
	}
	if len(rc) > 0 {
		return rc
	}
	return configs
}

// AddScraper adds a scraper.Scraper and checks if any steps were flagged,
// if flags are found any steps that have not been flagged are excluded from the scraper.Scraper.
// This exclusion prevents the step from being executed.
func (c *Config) AddScraper(id string, s *scraper.Scraper) {
	s.SetTag(attributes.NewTag(c.Id, id))
	c.Scrapers = append(c.Scrapers, s)
}
