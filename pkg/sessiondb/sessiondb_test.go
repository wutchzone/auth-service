package sessiondb

import (
	"testing"
	"time"
)

func TestNewSessionDB(t *testing.T) {
	// Create instance of redis
	d, err := NewSessionDB("localhost:7000", "", 0)
	//d.client.FlushAll()

	if err != nil {
		t.Errorf("Error creading redis connection expected nil, got %v", err)
		t.FailNow()
	}

	// Add record
	value := "testValue"
	key := "testKey"
	if err := d.setRecord(key, value, 0); err != nil {
		t.Errorf("Error setting record, expected nil, got: %v", err)
		t.FailNow()
	}

	if got, err := d.getRecord(key); err != nil || got != value {
		t.Errorf("Error getting error that was set. Expected: %v got %v with error %v", value, got, err)
	}

	// Remove record
	if err := d.removeRecord(key); err != nil {
		t.Errorf("Error removing key. Expected nil, got %v", err)
	}

	// Add timed record
	if err := d.setRecord(key, value, time.Second*2); err != nil {
		t.Errorf("Error adding key key. Expected nil, got %v", err)
	}

	// Check for timed record before expire
	if got, err := d.getRecord(key); err != nil || got != value {
		t.Errorf("Error getting error that was set. Expected: %v got %v with error %v", value, got, err)
	}

	// Check for timed record after expire
	time.Sleep(time.Second*2)
	if _, err := d.getRecord(key); err == nil {
		t.Errorf("Error getting error that was unset. Expected: error got nil")
	}
}
