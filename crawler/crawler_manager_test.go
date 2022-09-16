package crawler

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/database"
	"github.com/mmaaskant/gro-crop-scraper/test/httpserver"
	"net/http"
	"reflect"
	"testing"
)

var expected = map[string]*database.Entity{
	"http://localhost:8080/extract-1.html/": database.NewEntity(
		database.DbScrapedDataTableName,
		map[string]any{
			"_id":        nil,
			"tag":        "test",
			"origin":     "test_html",
			"url":        "http://localhost:8080/extract-1.html/",
			"html":       nil,
			"created_at": nil,
			"updated_at": nil,
		},
	),
	"http://localhost:8080/extract-2.html/": database.NewEntity(
		database.DbScrapedDataTableName,
		map[string]any{
			"_id":        nil,
			"tag":        "test",
			"origin":     "test_html",
			"url":        "http://localhost:8080/extract-2.html/",
			"html":       nil,
			"created_at": nil,
			"updated_at": nil,
		},
	),
}

func TestManager_Start(t *testing.T) {
	url := httpserver.StartTestHttpServer(t)
	db, err := database.NewDb(database.NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to database, error: %s", err)
	}
	m := NewCrawlerManager(db)
	c := NewHtmlCrawler("test_html", &http.Client{})
	c.tag = "test"
	c.AddDiscoveryUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?discovery-(\d*)(\.html)\/?`, url))
	c.AddExtractUrlRegex(fmt.Sprintf(`(https?:\/\/)?%s\/?extract-(\d*)(\.html)\/?`, url))
	m.RegisterCrawler(c, []*Call{NewCrawlerCall(
		fmt.Sprintf("http://%s/", url),
		DiscoverUrlType,
		http.MethodGet,
		nil,
		nil,
	)})
	m.Start(10)
	entities, err := db.GetMany(database.DbScrapedDataTableName, map[string]any{"tag": "test"})
	if err != nil {
		t.Errorf("Failed to fetch results from DB, error: %s", err)
	}
	for _, e := range entities {
		ex := expected[fmt.Sprint(e.Data["url"])]
		ex.Id = e.Id
		ex.Data["_id"] = e.Data["_id"]
		ex.Data["html"] = e.Data["html"]
		ex.CreatedAt = e.CreatedAt
		ex.Data["created_at"] = e.Data["created_at"]
		if e.Id == nil {
			t.Errorf("Entity %v does not have an ID.", e)
		}
		if e.Data["tag"] == nil {
			t.Errorf("Entity %v does not have a tag.", e)
		}
		if e.Data["origin"] == nil {
			t.Errorf("Entity %v does not have an origin.", e)
		}
		if e.CreatedAt == nil {
			t.Errorf("Entity %v does not have a created_at timestamp.", e)
		}
		if e.UpdatedAt != nil {
			t.Errorf("Entity %v has an updated_at timestamp, expected nil.", e)
		}
		if e.Data["html"] == nil {
			t.Errorf("Entity %v does not have html.", e)
		}
		if !reflect.DeepEqual(*e, *ex) {
			t.Errorf("Got entity %v, expected: %v", e, ex)
		}
	}
	err = db.DeleteMany(database.DbScrapedDataTableName, map[string]any{"tag": "test"})
	if err != nil {
		t.Errorf("Failed to tear down test data, error: %s", err)
	}
}
