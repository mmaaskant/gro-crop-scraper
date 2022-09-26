package scraper

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/filter"
	"log"
	"os"
	"strconv"
)

// Manager oversees all registered Scraper instances and its components.
type Manager struct {
	crawlerManager *crawler.Manager
	filterManager  *filter.Manager
	scrapers       []*Scraper
}

func NewManager(db *database.Db) *Manager {
	return &Manager{
		crawler.NewManager(db),
		filter.NewManager(db),
		make([]*Scraper, 0),
	}
}

func (m *Manager) RegisterScrapers(scrapers []*Scraper) {
	for _, s := range scrapers {
		m.RegisterScraper(s)
	}
}

func (m *Manager) RegisterScraper(s *Scraper) {
	if s.Crawler != nil && s.Calls != nil {
		m.crawlerManager.RegisterCrawler(s.Crawler, s.Calls)
	}
	if s.Filter != nil {
		m.filterManager.RegisterFilter(s.Filter)
	}
	m.scrapers = append(m.scrapers, s)
}

// Start starts Scraper and its components and waits till all components have finished running.
func (m *Manager) Start() {
	// TODO: Cache current process and restart where left off?
	// TODO: Only start if components have been found that manager is responsible for
	m.crawlerManager.Start(m.getWorkerCount("GOPHERVISOR_CRAWLER_WORKER_COUNT"))
	m.filterManager.Start(m.getWorkerCount("GOPHERVISOR_FILTER_WORKER_COUNT"))
}

// getWorkerCount gets a worker count from an env variable and attempts to convert it to an int.
func (m *Manager) getWorkerCount(env string) int {
	workerCount, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		log.Fatalf("Could not convert env variable %s with value %v to int",
			fmt.Sprintf("${%s}", env), os.Getenv(env),
		)
	}
	return workerCount
}
