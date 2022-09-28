package database

import "time"

// Driver holds functions to communicate with a database in a streamlined fashion,
// and is meant to be interchangeable so databases types can be easily switched if required.
// Driver currently does not support context.Context.
type Driver interface {
	connect() error
	GetOne(table string, params map[string]any) (*Entity, error)
	GetMany(table string, params map[string]any) (ResultIterator, error)
	InsertOne(e *Entity) error
	InsertMany(entities []*Entity) error
	UpdateOne(e *Entity) error
	UpdateMany(table string, filter map[string]any, update map[string]any) error
	DeleteOne(e *Entity) error
	DeleteMany(table string, filter map[string]any) error
}

type ResultIterator interface {
	Next() (*Entity, error)
}

// Entity holds results from DB queries and is used to interact with Driver.
// All data should be stored in Data, however it maps Id, CreatedAt and UpdatedAt for ease of access.
type Entity struct {
	Id        any
	Table     string
	Data      map[string]any
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func NewEntity(table string, data map[string]any) *Entity {
	return &Entity{
		Id:        nil,
		Table:     table,
		Data:      data,
		CreatedAt: nil,
		UpdatedAt: nil,
	}
}
