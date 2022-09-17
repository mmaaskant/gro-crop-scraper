package config

import (
	"flag"
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
)

// Config holds an origin tag that is passed on to all components to identify their origin
// and several maps holding components that can be run by scraper.Scraper.
type Config struct {
	Origin   string
	Crawlers map[crawler.Crawler][]*crawler.Call
}

// newConfig returns a new instance of Config.
func newConfig(origin string) *Config {
	return &Config{
		origin,
		make(map[crawler.Crawler][]*crawler.Call, 0),
	}
}

// AddCrawler adds an instance of Crawler and tags it with Config's origin.
func (c *Config) AddCrawler(cr crawler.Crawler, calls []*crawler.Call) {
	cr.SetOrigin(c.Origin)
	c.Crawlers[cr] = calls
}

// GetRegisteredConfigs returns all Config instances by default,
// if command origin flags are provided it will only return tagged Config instances.
func GetRegisteredConfigs() []*Config {
	registeredConfigs := make(map[*bool]*Config)
	configs := getConfigs()
	for _, c := range configs {
		f := flag.Bool(c.Origin, false, fmt.Sprintf("Add %s config to Scraper, runs all configs if none are added.", c.Origin))
		registeredConfigs[f] = c
	}
	flag.Parse()
	fc := getFlaggedConfigs(registeredConfigs)
	if len(fc) > 0 {
		return fc
	}
	return configs
}

// TODO: Automate this through reflection? Configs could be excluded through an annotation?
// getConfigs returns all Config instances that should and/or are ready to be executed.
func getConfigs() []*Config {
	return []*Config{
		NewBurpeeConfig(),
	}
}

// getFlaggedConfigs returns all instances of Config which have been flagged by command origin flags.
func getFlaggedConfigs(registeredConfigs map[*bool]*Config) []*Config {
	fc := make([]*Config, 0)
	for b, c := range registeredConfigs {
		if *b == true {
			fc = append(fc, c)
		}
	}
	return fc
}
