package accountdb

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
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
func (d *DB) getAccount(name string) (*User, error) {
	var u User

	r := d.db.Collection(d.defaultCollectionName).FindOne(nil, bson.M{"name": name})

	if err := r.Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

//SaveUser to the DB
func (d *DB) saveAccount(u User) error {
	if u, _ := d.getAccount(u.Username); u != nil {
		return errors.New("User already exists.")
	}

	if _, err := d.db.Collection(d.defaultCollectionName).InsertOne(nil, u); err != nil {
		return errors.New("Error while saving account to the DB.")
	} else{
		return nil
	}

}

func (d *DB) deleteAccount(u User) error {
	_, err := d.db.Collection(d.defaultCollectionName).DeleteOne(nil, bson.M{"name": u.Username})
	return err
}

// Modify user's data in the DB
func (d *DB) updateAccount(u User) error {
	_, err := d.db.Collection(d.defaultCollectionName).UpdateOne(nil, bson.M{"name": u.Username}, bson.M{"$set": u})
	return err
}
