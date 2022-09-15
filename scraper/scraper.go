package scraper

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"log"
	"os"
	"strconv"
)

// Scraper mediates between all components used for scraping, and oversees their process.
type Scraper struct {
	crawlerManager *crawler.Manager
}

// NewScraper returns a new instance of Scraper.
func NewScraper(db *database.Db) *Scraper {
	return &Scraper{
		crawler.NewCrawlerManager(db),
	}
}

// Start starts Scraper and its components and waits till all components have finished running.
func (s *Scraper) Start() {
	s.crawl()
}

// RegisterCrawler registers a new Crawler within crawler.CrawlerManager,
// provided crawler.Call instances will be crawled once scraper starts.
func (s *Scraper) RegisterCrawler(c crawler.Crawler, calls []*crawler.Call) {
	s.crawlerManager.RegisterCrawler(c, calls)
}

// crawl iterates over all registered crawlers and starts crawling their provided crawler.Call instances.
func (s *Scraper) crawl() {
	s.crawlerManager.Start(s.getWorkerCount("GOPHERVISOR_CRAWLER_WORKER_COUNT"))
}

// getWorkerCount gets a worker count from an env variable and attempts to convert it to an int.
func (s *Scraper) getWorkerCount(env string) int {
	workerCount, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		log.Fatalf("Could not convert env variable %s with value %v to int",
			fmt.Sprintf("${%s}", env), os.Getenv(env),
		)
	}
	return workerCount
}
