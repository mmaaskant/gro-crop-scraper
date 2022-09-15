package crawler

import (
	"github.com/mmaaskant/gophervisor/supervisor"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"log"
	"reflect"
)

// Manager oversees all registered Crawler instances.
type Manager struct {
	db       *database.Db
	crawlers map[Crawler][]*Call
}

// NewCrawlerManager returns a new instance of CrawlerManager.
func NewCrawlerManager(db *database.Db) *Manager {
	return &Manager{
		db,
		make(map[Crawler][]*Call),
	}
}

// crawlerJob holds a Crawler and a Call and is used to pass on units of work to Manager's workers.
type crawlerJob struct {
	c    Crawler
	call *Call
}

// newCrawlerJob returns a new instance of crawlerJob.
func newCrawlerJob(c Crawler, call *Call) *crawlerJob {
	return &crawlerJob{
		c:    c,
		call: call,
	}
}

// RegisterCrawler registers a new Crawler within crawler.CrawlerManager,
// provided crawler.Call instances will be crawled once scraper starts.
func (cm *Manager) RegisterCrawler(c Crawler, calls []*Call) {
	cm.crawlers[c] = calls
}

// Start begins crawling using the provided Crawler and Call instances,
// a supervisor.Supervisor instance is used to crawl concurrently.
func (cm *Manager) Start(amountOfWorkers int) {
	sv, p, _ := helper.StartSupervisor(amountOfWorkers, cm.crawl)
	for c, calls := range cm.crawlers {
		for _, call := range calls {
			p.Publish(newCrawlerJob(c, call))
		}
	}
	sv.Shutdown()
}

// crawl receives crawlerJob instances and handles them,
// this function is registered within supervisor.Supervisor as a worker.
func (cm *Manager) crawl(p *supervisor.Publisher, d any, rch chan any) {
	cj, ok := d.(*crawlerJob)
	if !ok {
		log.Fatalf("Scraper Crawl() expected instance of %s, got %s", "*crawlerJob", reflect.TypeOf(d))
	}
	cd := cj.c.Crawl(cj.call)
	for _, foundCall := range cd.FoundCalls {
		p.Publish(newCrawlerJob(cj.c, foundCall))
	}
	if cj.call.UrlType == ExtractUrlType {
		err := cm.db.InsertOne(database.NewEntity("scraped_html", map[string]any{
			"tag":  cd.Tag,
			"url":  cd.Call.Url,
			"html": cd.Data,
		}))
		if err != nil {
			log.Printf("Scraper failed to insert crawled HTML, error: %s", err)
		}
	}
}
