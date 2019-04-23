package sessiondb

import (
	"os"
	"testing"
	"time"
)

var REDIS_URL = "localhost:6379"

func TestNewSessionDB(t *testing.T) {
	// Try to get redis url from system env
	env := os.Getenv("REDIS")

	if env != "" {
		REDIS_URL = env
	}

	// Create instance of redis
	d, err := NewSessionDB(REDIS_URL, "", 0)
	//d.client.FlushAll()

	if err != nil {
		t.Errorf("Error creading redis connection expected nil, got %v", err)
		t.FailNow()
	}

	// Add record
	value := "testValue"
	key := "testKey"
	if err := d.SetRecord(key, value, 0); err != nil {
		t.Errorf("Error setting record, expected nil, got: %v", err)
		t.FailNow()
	}

	if got, err := d.GetRecord(key); err != nil || got != value {
		t.Errorf("Error getting error that was set. Expected: %v got %v with error %v", value, got, err)
	}

	// Remove record
	if err := d.RemoveRecord(key); err != nil {
		t.Errorf("Error removing key. Expected nil, got %v", err)
	}

	// Add timed record
	if err := d.SetRecord(key, value, time.Second*2); err != nil {
		t.Errorf("Error adding key key. Expected nil, got %v", err)
	}

	// Check for timed record before expire
	if got, err := d.GetRecord(key); err != nil || got != value {
		t.Errorf("Error getting error that was set. Expected: %v got %v with error %v", value, got, err)
	}

	// Check for timed record after expire
	time.Sleep(time.Second * 2)
	if _, err := d.GetRecord(key); err == nil {
		t.Errorf("Error getting error that was unset. Expected: error got nil")
	}
}
