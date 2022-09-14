package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"reflect"
	"time"
)

// MongoDbDriver communicates with an instance of MongoDB and offers queries through the Driver interface.
// MongoDbDriver is currently not concurrency proof and does not cache any results.
type MongoDbDriver struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDbDriver returns a new instance of MongoDbDriver.
func NewMongoDbDriver() *MongoDbDriver {
	return &MongoDbDriver{}
}

// Connect attempts to connect to an instance of MongoDB using the provided .env variables.
func (mdd *MongoDbDriver) Connect() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err == nil {
		mdd.client = client
		mdd.db = client.Database(os.Getenv("MONGO_INITDB_DATABASE"))
	}
	return err
}

// GetOne retrieves 1 row from MongoDB and maps it to an instance of Entity which it returns.
func (mdd *MongoDbDriver) GetOne(table string, params map[string]any) (*Entity, error) {
	var r bson.M
	err := mdd.db.Collection(table).FindOne(context.TODO(), mdd.bsonMarshal(params)).Decode(&r)
	if err != nil {
		return nil, err
	}
	e := mdd.newEntity(table, r)
	return e, err
}

// GetMany retrieves multiple rows from MongoDB and maps it to instances of Entity which it returns.
// An optional limit is available to fetch results in batches.
func (mdd *MongoDbDriver) GetMany(table string, params map[string]any, limit ...int) ([]*Entity, error) {
	l := int64(limit[0])
	opts := options.FindOptions{Limit: &l}
	c, err := mdd.db.Collection(table).Find(context.TODO(), mdd.bsonMarshal(params), &opts)
	if err != nil {
		return nil, err
	}
	var r []bson.D
	err = c.All(context.TODO(), &r)
	if err != nil {
		return nil, err
	}
	entities := make([]*Entity, 0)
	for _, v := range r {
		e := mdd.newEntity(table, v.Map())

		entities = append(entities, e)
	}
	return entities, err
}

// newEntity returns a new instance of Entity and sets values based on the retrieved MongoDB result set.
// This function is used to map results received from MongoDB to the Entity format.
func (mdd *MongoDbDriver) newEntity(table string, m bson.M) *Entity {
	e := NewEntity(table, m)
	e.Id = m["_id"]
	e.Data = m
	e.CreatedAt = mdd.convertToTime(m["created_at"])
	e.UpdatedAt = mdd.convertToTime(m["updated_at"])
	return e
}

// bsonMarshal encodes the given param to BSON.
func (mdd *MongoDbDriver) bsonMarshal(d any) []byte {
	r, err := bson.Marshal(d)
	if err != nil {
		log.Fatalf("BSON Marshal err: %s, Value: %v", err, d)
	}
	return r
}

// convertToTime converts the primitive.DateTime received from MongoDB to time.Time
func (mdd *MongoDbDriver) convertToTime(v any) *time.Time {
	pdt, ok := v.(primitive.DateTime)
	if !ok {
		log.Fatalf("Expected instance of primitive.DateTime while converting to *time.Time, got: %s", reflect.TypeOf(v))
	}
	t := pdt.Time()
	return &t
}

// InsertOne inserts 1 row into MongoDB based on the provided Entity.
func (mdd *MongoDbDriver) InsertOne(e *Entity) error {
	mdd.setCreatedAt(e)
	data := mdd.bsonMarshal(e.Data)
	ir, err := mdd.db.Collection(e.Table).InsertOne(context.TODO(), data)
	if err == nil {
		e.Id = ir.InsertedID
	}
	return err
}

// InsertMany inserts multiple rows in MongoDB based on the provided Entity instances.
func (mdd *MongoDbDriver) InsertMany(entities []*Entity) error {
	for table, m := range mdd.mapData(entities) {
		imr, err := mdd.db.Collection(table).InsertMany(context.TODO(), m)
		if err != nil {
			return err
		}
		for i, id := range imr.InsertedIDs {
			entities[i].Id = id
		}
	}
	return nil
}

// mapData iterates over the given Entity instances, converts them to BSON and maps them,
// this allows insertion of entities in bulk.
func (mdd *MongoDbDriver) mapData(entities []*Entity) map[string][]any {
	m := make(map[string][]any)
	for _, e := range entities {
		mdd.setCreatedAt(e)
		data := mdd.bsonMarshal(e.Data)
		m[e.Table] = append(m[e.Table], data)
	}
	return m
}

// setCreatedAt sets the Entity.CreatedAt value to time.Now() and prepares Entity.UpdatedAt value for future use.
func (mdd *MongoDbDriver) setCreatedAt(e *Entity) {
	t := time.Now()
	e.CreatedAt = &t
	e.Data["created_at"] = time.Now()
	e.Data["updated_at"] = nil
}

// UpdateOne updates 1 instance of Entity.
func (mdd *MongoDbDriver) UpdateOne(e *Entity) error {
	_, err := mdd.db.Collection(e.Table).UpdateByID(context.TODO(), e.Id, mdd.bsonMarshal(e.Data))
	return err
}

// UpdateMany updates multiple rows within MongoDB based on the provided filter and fields to be updated.
func (mdd *MongoDbDriver) UpdateMany(table string, filter map[string]any, update map[string]any) error {
	_, err := mdd.db.Collection(table).UpdateMany(context.TODO(), mdd.bsonMarshal(filter), mdd.bsonMarshal(update))
	return err
}

// setUpdatedAt sets Entity.UpdatedAt to time.Now().
func (mdd *MongoDbDriver) setUpdatedAt(e *Entity) {
	t := time.Now()
	e.UpdatedAt = &t
	e.Data["updated_at"] = time.Now()
}

// DeleteOne deletes 1 row within MongoDB based on the provided Entity.
func (mdd *MongoDbDriver) DeleteOne(e *Entity) error {
	_, err := mdd.db.Collection(e.Table).DeleteOne(context.TODO(), bson.D{{"_id", e.Id}})
	return err
}

// DeleteMany deletes multiple rows within MongoDB based on the provided table and filter.
func (mdd *MongoDbDriver) DeleteMany(table string, filter map[string]any) error {
	_, err := mdd.db.Collection(table).DeleteMany(context.TODO(), mdd.bsonMarshal(filter))
	return err
}
