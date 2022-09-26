package filter

import (
	"fmt"
	"github.com/mmaaskant/gophervisor/supervisor"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"log"
	"reflect"
)

// TODO: Add comments

type Manager struct {
	db      *database.Db
	filters []Filter
}

func NewManager(db *database.Db) *Manager {
	return &Manager{
		db,
		make([]Filter, 0),
	}
}

type filterJob struct {
	filter Filter
	entity *database.Entity
}

func newFilterJob(f Filter, e *database.Entity) *filterJob {
	return &filterJob{
		f.Clone(),
		e,
	}
}

func (m *Manager) RegisterFilters(filters []Filter) {
	for _, f := range filters {
		m.RegisterFilter(f)
	}
}

func (m *Manager) RegisterFilter(f Filter) {
	m.filters = append(m.filters, f)
}

func (m *Manager) Start(amountOfWorkers int) {
	sv, p, _ := helper.StartSupervisor(amountOfWorkers, m.filter)
	for _, f := range m.filters {
		getMany := func() ([]*database.Entity, error) {
			return m.db.GetMany(database.ScrapedDataTableName, map[string]any{"scraper_id": f.GetScraperId()}, 50)
		}
		for entities, err := getMany(); len(entities) > 0; entities, err = getMany() {
			if err != nil {
				log.Printf("Failed to fetch entities from DB table %s, error: %s", database.ScrapedDataTableName, err)
			}
			for _, e := range entities {
				p.Publish(newFilterJob(f, e))
				if err = m.db.DeleteOne(e); err != nil { // TODO: Entity is inaccessible too early in case it is needed in the future
					log.Printf("Failed to delete entity %v, error: %s", e, err)
				}
			}
		}
	}
	sv.Shutdown()
}

func (m *Manager) filter(p *supervisor.Publisher, d any, rch chan any) {
	var fj *filterJob
	fj, ok := d.(*filterJob)
	if !ok {
		log.Fatalf("Expected instance of %s, got %s", reflect.TypeOf(fj), reflect.TypeOf(d))
	}
	if data := fj.filter.Filter(fmt.Sprint(fj.entity.Data["data"])); data != nil && len(data) > 0 {
		// TODO: Insert filtered data
		fmt.Println(data)
	}
}
