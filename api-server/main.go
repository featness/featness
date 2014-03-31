package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/tsuru/config"
	"log"
	"os"
)

const version = "0.1.0"

func parseFlags(args []string) (*string, *bool) {
	flagSet := flag.NewFlagSet("configuration", flag.ExitOnError)
	configFile := flagSet.String("config", "/etc/featness-api.conf", "Featness API configuration file")
	gVersion := flagSet.Bool("version", false, "Print version and exit")
	flagSet.Parse(args)

	return configFile, gVersion
}

func loadConfigFile(path *string) {
	err := config.ReadAndWatchConfigFile(*path)
	if err != nil {
		msg := `Could not find featness-api config file. Searched on %s.
	For an example conf check featness-api/etc/featness-api.conf file.\n %s`
		log.Panicf(msg, *path, err)
	}
}

func getHandlers() *pat.Router {
	return &pat.Router{}
}

func main() {
	configFile, gVersion := parseFlags(os.Args)

	if *gVersion {
		fmt.Printf("featness-api version %s\n", version)
		return
	}

	loadConfigFile(configFile)

	//router := getHandlers()
}
