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

func getRouter() *pat.Router {
	router := pat.New()
	router.Get("/healthcheck", api.Healthcheck)

	return router
}

func main() {
	fmt.Printf("%s", os.Args[1:])
	configFile, gVersion := parseFlags(os.Args[1:])

	if gVersion {
		fmt.Printf("featness-api version %s\n", api.Version)
		return
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	loadConfigFile(configFile, logger)

	router := getRouter()
	log.Println("featness-api running at http://localhost:8000...")
	http.ListenAndServe(":8000", router)
}
