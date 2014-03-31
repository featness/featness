package main

import (
	"flag"
	"fmt"
	"github.com/globoi/featness/api"
	"github.com/gorilla/pat"
	"net/http"
	"os"
)

func parseFlags(args []string) bool {
	flagSet := flag.NewFlagSet("configuration", flag.ExitOnError)
	gVersion := flagSet.Bool("version", false, "Print version and exit")
	flagSet.Parse(args)

	return *gVersion
}

func serveAngular(w http.ResponseWriter, r *http.Request) {
	data, _ := Asset("dashboard/index.html")
	fmt.Printf("PASSOU AKI\n")
	w.Write(data)
}

func serveScripts(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get(":path")
	data, _ := Asset(fmt.Sprintf("dashboard/scripts/%s", path))
	w.Write(data)
}

func serveStyles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get(":path")
	data, _ := Asset(fmt.Sprintf("dashboard/styles/%s", path))
	w.Write(data)
}

func getRouter() *pat.Router {
	router := pat.New()
	router.Get("/scripts/{path:.+}", serveScripts)
	router.Get("/styles/{path:.+}", serveStyles)
	router.Get("/", serveAngular)

	return router
}

func main() {
	gVersion := parseFlags(os.Args[1:])

	if gVersion {
		fmt.Printf("featness-dashboard version %s\n", api.Version)
		return
	}

	router := getRouter()
	http.ListenAndServe(":8080", router)
}
