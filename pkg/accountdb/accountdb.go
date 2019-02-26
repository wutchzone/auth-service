package accountdb

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBItem interface {
	Name() string
}

type DB struct {
	client                *mongo.Client
	db                    *mongo.Database
	defaultCollectionName string
}

// NewSession returns pointer to the new sessiondb
func NewAccountDBConnection(addr string, table string, defaultCollection string) (*DB, error) {
	d := &DB{}

	// Create client instance
	c, err := mongo.NewClient(addr)
	if err != nil {
		return nil, err
	}

	// Try to connect to the DB
	err = c.Connect(context.Background())
	err = c.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// assign DB and client reference for future purpose
	d.defaultCollectionName = defaultCollection
	d.db = c.Database(table)
	d.client = c
	return d, nil
}

// GetUser from DB
func (d *DB) GetAccount(name string) *mongo.SingleResult {
	return d.db.Collection(d.defaultCollectionName).FindOne(nil, bson.M{"name": name})
}

//SaveUser to the DB
func (d *DB) SaveAccount(u DBItem) error {
	if _, err := d.db.Collection(d.defaultCollectionName).InsertOne(nil, u); err != nil {
		return errors.New("Error while saving account to the DB.")
	} else {
		return nil
	}

}

func (d *DB) GetAll() *mongo.Cursor {
	cursor, _ := d.db.Collection(d.defaultCollectionName).Find(nil, bson.M{})

	return cursor
}

func (d *DB) DeleteAccount(u string) error {
	_, err := d.db.Collection(d.defaultCollectionName).DeleteOne(nil, bson.M{"name": u})
	return err
}

// Modify user's data in the DB
func (d *DB) UpdateAccount(u DBItem) error {
	_, err := d.db.Collection(d.defaultCollectionName).UpdateOne(nil, bson.M{"name": u.Name()}, bson.M{"$set": u})
	return err
}
