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

func NewManager(db *database.Db) *Manager {
	return &Manager{
		db,
		make(map[Crawler][]*Call),
	}
}

// crawlerJob holds a Crawler and a Call and is used to pass on units of work to Manager's workers.
type crawlerJob struct {
	crawler Crawler
	call    *Call
}

func newCrawlerJob(c Crawler, call *Call) *crawlerJob {
	return &crawlerJob{
		crawler: c,
		call:    call,
	}
}

func (m *Manager) RegisterCrawlers(crawlers map[Crawler][]*Call) {
	for c, calls := range crawlers {
		m.RegisterCrawler(c, calls)
	}
}

func (m *Manager) RegisterCrawler(c Crawler, calls []*Call) {
	m.crawlers[c] = calls
}

// Start begins crawling using the provided Crawler and Call instances,
// a supervisor.Supervisor instance is used to crawl concurrently.
func (m *Manager) Start(amountOfWorkers int) {
	sv, p, _ := helper.StartSupervisor(amountOfWorkers, m.crawl)
	for c, calls := range m.crawlers {
		for _, call := range calls {
			p.Publish(newCrawlerJob(c, call))
		}
	}
	sv.Shutdown()
}

// crawl receives crawlerJob instances and handles them,
// this function is registered within supervisor.Supervisor as a worker.
func (m *Manager) crawl(p *supervisor.Publisher, d any, rch chan any) {
	var cj *crawlerJob
	cj, ok := d.(*crawlerJob)
	if !ok {
		log.Fatalf("Expected instance of %s, got %s", reflect.TypeOf(cj), reflect.TypeOf(d))
	}
	cd := cj.crawler.Crawl(cj.call)
	if cd.Error != nil {
		log.Printf("Crawler data contains error %s, skipping ...", cd.Error)
		return
	}
	for _, foundCall := range cd.FoundCalls {
		p.Publish(newCrawlerJob(cj.crawler, foundCall))
	}
	if cj.call.RequestType == ExtractRequestType {
		err := m.db.InsertOne(database.NewEntity(database.ScrapedDataTableName, map[string]any{
			"config_id":  cd.GetConfigId(),
			"scraper_id": cd.GetScraperId(),
			"url":        cd.Call.Request.URL.String(),
			"data":       cd.Data,
		}))
		if err != nil {
			log.Printf("Scraper failed to insert crawled HTML, error: %s", err)
		}
	}
}
