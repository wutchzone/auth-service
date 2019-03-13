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
	ServiceDB *accountdb.DB
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

	// Start up configuration
	if !Config.Dev.IgnoreStartup {
		sud, err := accountdb.NewAccountDBConnection(Config.GeneralDB.URL, Config.GeneralDB.Table, "general")
		if err != nil {
			panic(err)
		}
		r := sud.GetAccount("startup")

		if err != nil {
			panic(err)
		}

		i := &configuration.StartupConfiguration{}

		// Try if document exists
		if err := r.Decode(i); err != nil {
			ns := configuration.NewStartup()
			ns.FirstBoot = false
			sud.SaveAccount(ns)
			// First init
			// Push default service to the DB
			if Config.User.ServiceToken != "" {
				ServiceDB.SaveAccount(&accountdb.Service{
					ID:    Config.User.ServiceToken,
					Level: 999999999,
				})
			}

			fmt.Println("First initialization")
		} else {
			// Late init
		}
	}

	// Load services do ServiceDB
	SessionDB.Client.FlushAll()
	rslt := ServiceDB.GetAll()
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
