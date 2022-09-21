package database

import (
	"testing"
)

func TestMongoDbDriver_One(t *testing.T) {
	db := newDb(t)
	e := newScrapedHtmlEntity("https://example.com")
	err := db.InsertOne(e)
	if err != nil {
		t.Errorf("Failed to insert single row into DB, err: %s", err)
	}
	e.Data["updated"] = true
	err = db.UpdateOne(e)
	if err != nil {
		t.Errorf("Failed to update single row in DB, error: %s", err)
	}
	_, err = db.GetOne(ScrapedDataTableName, map[string]any{"_id": e.Id})
	if err != nil {
		t.Errorf("Failed to get single row from DB, error: %s", err)
	}
	if e.UpdatedAt == nil || e.Data["updated"] != true {
		if e.UpdatedAt == nil {
			t.Errorf("Entity %v UpdatedAt is nil, expected timestamp.", e)
		}
		if e.Data["updated"] != true {
			t.Errorf("Entity %v data.updated is %v, expected: %v", e, e.Data["updated"], true)
		}
	}
	err = db.DeleteOne(e)
	if err != nil {
		t.Errorf("Failed to delete single row in DB, error: %s", err)
	}
}

func TestMongoDbDriver_Many(t *testing.T) {
	db := newDb(t)
	err := db.InsertMany([]*Entity{
		newScrapedHtmlEntity("https://example.com/1/"),
		newScrapedHtmlEntity("https://example.com/2/"),
		newScrapedHtmlEntity("https://example.com/3/"),
	})
	if err != nil {
		t.Errorf("Failed to insert multiple rows into DB, err: %s", err)
	}
	if err != nil {
		t.Errorf("Failed to get multiple rows from DB, error: %s", err)
	}
	err = db.UpdateMany(ScrapedDataTableName, map[string]any{"updated": false}, map[string]any{"updated": true})
	if err != nil {
		t.Errorf("Failed to update multiple rows in DB, error: %s", err)
	}
	entities, err := db.GetMany(ScrapedDataTableName, map[string]any{"config_id": "test"})
	for _, entity := range entities {
		if entity.UpdatedAt == nil {
			t.Errorf("Entity %v UpdatedAt is nil, expected timestamp.", entity)
		}
		if entity.Data["updated"] != true {
			t.Errorf("Entity %v data.updated is %v, expected: %v", entity, entity.Data["updated"], true)
		}
	}
	err = db.DeleteMany(ScrapedDataTableName, map[string]any{"config_id": "test"})
	if err != nil {
		t.Errorf("Failed to delete multiple rows in DB, error: %s", err)
	}
}

func newDb(t *testing.T) *Db {
	db, err := NewDb(NewMongoDbDriver())
	if err != nil {
		t.Errorf("Failed to connect to DB using MongoDriver, error: %s", err)
	}
	return db
}

func newScrapedHtmlEntity(url string) *Entity {
	return NewEntity(
		ScrapedDataTableName,
		map[string]any{"url": url, "config_id": "test", "scraper_id": "test_html", "updated": false},
	)
}
