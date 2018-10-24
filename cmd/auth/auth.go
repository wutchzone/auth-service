package main

import (
	"fmt"
	"net/http"

	"github.com/wutchzone/auth-service/api"
	"github.com/wutchzone/auth-service/pkg/session"
	"github.com/wutchzone/auth-service/pkg/userdb"
)

// Config is reference for service configuration
var Config *parser.ConfigJSON

func init() {
	// f, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
	// if err != nil {
	// 	fmt.Println("Can not read file %n", f)
	// }
	// json.NewDecoder(strings.NewReader(string(f))).Decode(&Config)
	// fmt.Println(Config.API)
	fmt.Println(userdb.NewSession("localhost:7000", ""))
	fmt.Println(session.NewSession("localhost:7000", ""))

}

func main() {
	r := InitRoutes()

	http.ListenAndServe(":7080", r)
}
