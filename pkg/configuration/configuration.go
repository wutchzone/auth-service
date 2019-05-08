package configuration

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// Dev definition
type Dev struct {
	IgnoreStartup bool `json:"ignore_startup"`
}

// Route definition
type Route struct {
	From  string
	To    string
	Level int
}

// Role definition
type Role struct {
	Name  string // Assigned name
	Level int    // Access level
}

type DB struct {
	URL      string // URL to the connected DB
	Name     string // Name for connecting to the DB
	Password string // Password for connection to the DB
	Table    string // Used table name
}

// DefaultUser is created when you firstly start your app. It has admin privilegies.
// It is recommended to chanfe it to the different password.
type DefaultUser struct {
	Name         string "admin"
	Password     string "admin"
	ServiceToken string // Default token that is immediately loaded to the DB after first boot
}

type Configuration struct {
	DefaultPort int         `json:"default_port"` // Port where service will be listening
	Roles       []Role      `json:"roles"`        // Roles definition
	SessionDB   DB          `json:"sessiondb"`    // DB for storing sessiondb information (Only Redis is currently supported)
	AccessDB    DB          `json:"accessdb"`     // Parameters for connecting to the DB where user data's are stored (Only MongoDB is currently supported)
	User        DefaultUser `json:"user"`         // Default admin account
	Dev         Dev         `json:"dev"`          // Configuration for development purpose
}

func NewConfiguration(path string) (*Configuration, error) {
	// Read config file
	f, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		return nil, errors.New("configuration file was not specified")
	}

	config := &Configuration{}
	// Parse configuration file
	if err := json.NewDecoder(strings.NewReader(string(f))).Decode(&config); err != nil {
		return nil, errors.New("configuration file is badly formatted")
	}

	return config, nil
}
