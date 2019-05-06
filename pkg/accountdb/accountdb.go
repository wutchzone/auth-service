package accountdb

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Singleton *DB

type DB struct {
	client *mongo.Client
	db     *mongo.Database
	config AccountConfiguration
}

type AccountConfiguration struct {
	AccoutCollectionName  string
	ServiceCollectionName string
	Address               string
}

// GetInstance returns pointer to the actual connection. Func expect address, name of the table and default
// collection name
func GetInstance(conf AccountConfiguration) *DB {
	if Singleton == nil {
		d := &DB{}
		d.config = conf
		// Create client instance
		c, err := mongo.NewClient(conf.Address)
		if err != nil {
			return nil
		}

		// Try to connect to the DB
		err = c.Connect(context.Background())
		err = c.Ping(context.Background(), nil)
		if err != nil {
			return nil
		}

		// assign DB and client reference for future purpose
		d.db = c.Database("auth")
		d.client = c

		Singleton = d
	}

	return Singleton
}

// GetUser from DB by specific name
func (d *DB) GetAccount(name string) (*User, error) {
	o := &User{}
	result := d.db.Collection(d.config.AccoutCollectionName).FindOne(nil, bson.M{"name": name})

	if err := result.Decode(o); err != nil {
		return nil, err
	}

	return o, nil
}

//SaveUser saves new user instance to the DB
func (d *DB) SaveUser(u User) error {
	if _, err := d.GetAccount(u.Name()); err == nil { // error == nil means that user was found succesfully
		return errors.New("user already exists")
	}

	if _, err := d.db.Collection(d.config.AccoutCollectionName).InsertOne(nil, u); err != nil {
		return errors.New("error while saving account to the db")
	} else {
		return nil
	}

}

func (d *DB) GetAllUsers() ([]User, error) {
	o := []User{}
	cur, err := d.db.Collection(d.config.AccoutCollectionName).Find(nil, bson.M{})

	if err != nil {
		return nil, errors.New("unable to retrieve all users")
	}

	for cur.Next(context.Background()) {
		var elem User

		err := cur.Decode(&elem)
		if err != nil {
			return nil, errors.New("unable to retrieve all users")
		}
		o = append(o, elem)
	}

	return o, nil
}

func (d *DB) DeleteUser(u string) error {
	_, err := d.db.Collection(d.config.AccoutCollectionName).DeleteOne(nil, bson.M{"name": u})
	return err
}

// Modify user's data in the DB
func (d *DB) UpdateUser(u User) error {
	_, err := d.db.Collection(d.config.AccoutCollectionName).UpdateOne(nil, bson.M{"name": u.Name()}, bson.M{"$set": u})
	return err
}
