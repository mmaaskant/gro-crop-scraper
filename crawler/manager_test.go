package crawler

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/attribute"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/test/httpserver"
	"net/http"
	"reflect"
	"testing"
)

var expected = map[string]*database.Entity{
	"http://localhost:8080/extract-1.html/": database.NewEntity(
		database.ScrapedDataTableName,
		map[string]any{
			"_id":        nil,
			"config_id":  "test",
			"scraper_id": "test_html",
			"url":        "http://localhost:8080/extract-1.html/",
			"data":       nil,
			"created_at": nil,
			"updated_at": nil,
		},
	),
	"http://localhost:8080/extract-2.html/": database.NewEntity(
		database.ScrapedDataTableName,
		map[string]any{
			"_id":        nil,
			"config_id":  "test",
			"scraper_id": "test_html",
			"url":        "http://localhost:8080/extract-2.html/",
			"data":       nil,
			"created_at": nil,
			"updated_at": nil,
		},
	),
}

func TestCrawlerManager_Start(t *testing.T) {
	url := httpserver.StartTestHttpServer(t)
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to database, error: %s", err)
	}
	m := NewManager(db)
	c := NewHtmlCrawler(&http.Client{})
	c.SetTag(attribute.NewTag("test", "test_html"))
	c.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	c.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	m.RegisterCrawler(c, []*Call{NewCall(
		NewRequest(http.MethodGet, fmt.Sprintf("http://%s/", url), nil),
		DiscoverRequestType,
	)})
	m.Start(10)
	iterator, err := db.GetMany(database.ScrapedDataTableName, map[string]any{"config_id": "test"})
	if err != nil {
		t.Errorf("Failed to initialise iterator, error: %s", err)
	}
	for e, _ := iterator.Next(); e != nil; e, _ = iterator.Next() {
		ex := expected[fmt.Sprint(e.Data["url"])]
		ex.Id = e.Id
		ex.Data["_id"] = e.Data["_id"]
		ex.Data["data"] = e.Data["data"]
		ex.CreatedAt = e.CreatedAt
		ex.Data["created_at"] = e.Data["created_at"]
		if e.Id == nil {
			t.Errorf("Entity %v does not have an ID.", e)
		}
		if e.Data["config_id"] == nil {
			t.Errorf("Entity %v does not have a config id.", e)
		}
		if e.Data["scraper_id"] == nil {
			t.Errorf("Entity %v does not have a scraper id.", e)
		}
		if e.CreatedAt == nil {
			t.Errorf("Entity %v does not have a created_at timestamp.", e)
		}
		if e.UpdatedAt != nil {
			t.Errorf("Entity %v has an updated_at timestamp, expected nil.", e)
		}
		if e.Data["data"] == nil {
			t.Errorf("Entity %v does not have data.", e)
		}
		if !reflect.DeepEqual(*e, *ex) {
			t.Errorf("Got entity %v, expected: %v", e, ex)
		}
	}
	err = db.DeleteMany(database.ScrapedDataTableName, map[string]any{"config_id": "test"})
	if err != nil {
		t.Errorf("Failed to tear down test data, error: %s", err)
	}
}
