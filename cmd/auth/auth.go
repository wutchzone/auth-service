package main

import (
	"encoding/json"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"github.com/wutchzone/auth-service/pkg/configuration"
	"github.com/wutchzone/auth-service/pkg/sessiondb"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Config is reference for service configuration
var (
	Config    *configuration.Configuration
	UserDB    *accountdb.DB
	ServiceDB *accountdb.DB
	SessionDB *sessiondb.SessionDB
)

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
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: Configuration loaded succesfully"))
	}

	// Init account DB
	if udb, err := accountdb.NewAccountDBConnection(Config.AccountDB.URL, Config.AccountDB.Table, "users"); err != nil {
		panic("Error connecting to the account DB")
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: AccountDB loaded succesfully"))
		UserDB = udb
	}

	// Init service DB
	if sdb, err := accountdb.NewAccountDBConnection(Config.ServiceDB.URL, Config.ServiceDB.Table, "services"); err != nil {
		panic("Error connecting to the service DB")
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: ServiceDB loaded succesfully"))
		ServiceDB = sdb
	}

	// Init sessiondb DB
	if sdb, err := sessiondb.NewSessionDB(Config.SessionDB.URL, Config.SessionDB.Password, 1); err != nil {
		panic("Error connecting to the sessiondb DB")
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: SessionDB loaded succesfully"))
		SessionDB = sdb
	}

	// Init SMTP
	fmt.Println("SMTP not implemented")

	// Load services do ServiceDB
	SessionDB.Client.FlushAll()
	rslt := ServiceDB.GetAll()
	ser, err := decodeServices(rslt)
	for _, i := range ser{
		SessionDB.SetRecord(i.Name(), "service", 0)
	}
	fmt.Println(emoji.Sprint(":white_check_mark: Services loaded succesfully"))

	fmt.Println("Everything started!")
}

func main() {
	r := InitRoutes()

	http.ListenAndServe(":7080", r)
}
