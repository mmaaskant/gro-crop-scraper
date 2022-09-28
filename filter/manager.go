package filter

import (
	"fmt"
	"github.com/mmaaskant/gophervisor/supervisor"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"go.mongodb.org/mongo-driver/mongo"
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
		iterator, err := m.db.GetMany(database.ScrapedDataTableName, map[string]any{"scraper_id": f.GetScraperId()})
		if err != nil {
			log.Panicf("Failed to initialise iterator, error: %s", err)
		}
		for e, _ := iterator.Next(); e != nil; e, _ = iterator.Next() {
			p.Publish(newFilterJob(f, e))
		}
	}
	sv.Shutdown()
}

func (m *Manager) filter(p *supervisor.Publisher, d any, rch chan any) {
	var fj *filterJob
	fj, ok := d.(*filterJob)
	if !ok {
		log.Panicf("Expected instance of %s, got %s", reflect.TypeOf(fj), reflect.TypeOf(d))
	}
	if data := fj.filter.Filter(fmt.Sprint(fj.entity.Data["data"])); data != nil {
		if fe, err := m.db.GetOne(database.FilteredDataTableName, map[string]any{"url": fj.entity.Data["url"]}); err != mongo.ErrNoDocuments {
			fe.Data["data"] = data
			if err = m.db.UpdateOne(fe); err != nil {
				log.Panicf("Failed to update filtered data, error: %s", err)
			}
		} else {
			e := database.NewEntity(
				database.FilteredDataTableName,
				map[string]any{
					"url":        fj.entity.Data["url"],
					"config_id":  fj.filter.GetConfigId(),
					"scraper_id": fj.filter.GetScraperId(),
					"data":       data,
				},
			)
			if err = m.db.InsertOne(e); err != nil {
				log.Panicf("Failed to insert filtered data, error: %s", err)
			}
			if err = m.db.DeleteOne(fj.entity); err != nil {
				log.Panicf("Failed to delete scraped data, error: %s", err)
			}
		}
	}
}
