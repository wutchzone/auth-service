package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/wutchzone/auth-service/api"
	"github.com/wutchzone/auth-service/pkg/session"
	"github.com/wutchzone/auth-service/pkg/userdb"
)

// Config is reference for service configuration
var Config *parser.ConfigJSON

func init() {
	f, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println("Can not read file", f)
	}
	json.NewDecoder(strings.NewReader(string(f))).Decode(&Config)

	errUserDB := userdb.NewSession(Config.Userdb.Address, Config.Userdb.Password)
	if errUserDB != nil {
		fmt.Println("Error connecting to the UserDB: ", errUserDB)
		return
	}

	errSessionDB := session.NewSession(Config.Sessiondb.Address, Config.Sessiondb.Password)
	if errSessionDB != nil {
		fmt.Println("Error connecting to the SessionDB: ", errSessionDB)
		return
	}

	fmt.Println("Everything started succesfully.")
}

func main() {
	r := InitRoutes()
	fmt.Println(":" + string(Config.Port))
	http.ListenAndServe(":7080", r)
}
