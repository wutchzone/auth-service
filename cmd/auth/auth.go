package main

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"github.com/wutchzone/auth-service/pkg/configuration"
	"github.com/wutchzone/auth-service/pkg/sessiondb"
	"net/http"
	"os"
)

// Config is reference for service configuration
var (
	Config    *configuration.Configuration
	UserDB    *accountdb.DB
	SessionDB *sessiondb.SessionDB
)

func init() {
	// Check if config file was passed
	if len(os.Args) == 1 || len(os.Args) > 2 {
		panic("Configuration file was not specified as an argument")
	}

	// Parse configuration file
	if conf, err := configuration.NewConfiguration(os.Args[1]); err != nil {
		panic("Configuration file is badly formatted")
	} else {
		Config = conf
		fmt.Println(emoji.Sprint(":white_check_mark: Configuration loaded succesfully"))
	}

	// Init account DB
	if udb := accountdb.GetInstance(accountdb.AccountConfiguration{
		AccoutCollectionName:  "users",
		ServiceCollectionName: "services",
		Address:               Config.SessionDB.URL,
	}); udb == nil {
		panic("Error connecting to the account & service DB")
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: AccountDB & ServiceDB loaded succesfully"))
		UserDB = udb
	}

	// Init sessiondb DB
	if sdb := sessiondb.GetInstance(sessiondb.SessionDBConfiguration{
		Address: Config.SessionDB.Table,
		TableID: 1,
	}); sdb == nil {
		panic("Error connecting to the sessiondb DB")
	} else {
		fmt.Println(emoji.Sprint(":white_check_mark: SessionDB loaded succesfully"))
		SessionDB = sdb
	}

	// Init SMTP
	fmt.Println("SMTP not implemented")

	ser, _ := decodeServices(rslt)
	for _, i := range ser {
		SessionDB.SetRecord(i.Name(), "service", 0)
	}
	fmt.Println(emoji.Sprint(":white_check_mark: Services loaded succesfully"))

	fmt.Println("Everything started!")
}

func main() {
	r := InitRoutes()

	http.ListenAndServe(":7080", r)
}
