package main

import (
	"encoding/json"
	"fmt"
	"github.com/wutchzone/auth-service/pkg/configuration"
	"github.com/wutchzone/auth-service/pkg/sessiondb"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Config is reference for service configuration
var Config *configuration.Configuration

func init() {
	// Check if config file was passed
	if len(os.Args) == 1 {
		panic("Configuration file was not specified as an argument")
	}

	// Read config file
	f, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		panic("Can not read configuration file")
	}

	// Parse configuration file
	if err := json.NewDecoder(strings.NewReader(string(f))).Decode(&Config); err != nil {
		panic("Configuration file is badly formatted")
	}

	// Init account DB
	if err := accountdb.NewSession(Config.AccountDB.URL, Config.Userdb.Password); err != nil {
		panic("Error connecting to the account DB")
	}

	// Init sessiondb DB
	if err := sessiondb.NewSession(Config.Sessiondb.Address, Config.Sessiondb.Password); err != nil {
		panic("Error connecting to the sessiondb DB")
	}

	fmt.Println("Everything started!")
}

func main() {
	r := InitRoutes()
	fmt.Println(":" + string(Config.Port))
	http.ListenAndServe(":7080", r)
}
