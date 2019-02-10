package accountdb

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
)

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
func (d *DB) GetAccount(name string) (*User, error) {
	var u User

	r := d.db.Collection(d.defaultCollectionName).FindOne(nil, bson.M{"name": name})

	if err := r.Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

//SaveUser to the DB
func (d *DB) SaveAccount(u User) error {
	return errors.New("")
}

func (d *DB) DeleteAccount(u User) error {
	return errors.New("")
}

func (d *DB) updateAccount(u User) error {
	return errors.New("")
}
