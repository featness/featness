package main

import (
	"flag"
	"fmt"
	"github.com/globoi/featness/api"
	"github.com/gorilla/pat"
	"github.com/tsuru/config"
	"log"
	"net/http"
	"os"
)

type Logger interface {
	Panicf(format string, v ...interface{})
}

func parseFlags(args []string) (string, bool) {
	flagSet := flag.NewFlagSet("configuration", flag.ExitOnError)
	configFile := flagSet.String("config", "/etc/featness-api.conf", "Featness API configuration file")
	gVersion := flagSet.Bool("version", false, "Print version and exit")
	flagSet.Parse(args)

	return *configFile, *gVersion
}

func loadConfigFile(path string, logger Logger) {
	err := config.ReadAndWatchConfigFile(path)
	if err != nil {
		msg := `Could not find featness-api config file. Searched on %s.
	For an example conf check api/etc/local.conf file.\n %s`
		logger.Panicf(msg, path, err)
	}
}

func crossDomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-Auth-Data, X-Auth-Token")
	w.WriteHeader(http.StatusOK)
}

func AllowCrossDomainFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "Accept, Content-Type, X-Auth-Data, X-Auth-Token")
		handler(w, r)
	}
}

func getRouter() *pat.Router {
	router := pat.New()
	router.Get("/healthcheck", AllowCrossDomainFunc(api.Healthcheck))
	router.Post("/authenticate/google", AllowCrossDomainFunc(api.AuthenticateWithGoogle))
	router.Post("/authenticate/facebook", AllowCrossDomainFunc(api.AuthenticateWithFacebook))
	router.Add("OPTIONS", "/", http.HandlerFunc(crossDomain))

	return router
}

func connectToMongo() error {
	hosts, err := config.GetString("mongo:hosts")
	if err != nil {
		return fmt.Errorf("Could not find MongoDB host information (%s).", err)
	}

	database, err := config.GetString("mongo:database")
	if err != nil {
		return fmt.Errorf("Could not find MongoDB database information (%s).", err)
	}

	username, _ := config.GetString("mongo:username")
	password, _ := config.GetString("mongo:password")

	fmt.Println("Connecting to mongo at ", hosts, " database: ", database, " user: ", username, " password: ", password)

	return api.MongoStartup("featness", hosts, database, username, password)
}

func main() {
	configFile, gVersion := parseFlags(os.Args[1:])

	if gVersion {
		fmt.Printf("featness-api version %s\n", api.Version)
		return
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	loadConfigFile(configFile, logger)
	err := connectToMongo()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := getRouter()
	log.Println("featness-api running at http://localhost:8000...")
	http.ListenAndServe(":8000", router)
}
