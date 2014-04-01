package main

import (
	"flag"
	"fmt"
	"github.com/globoi/featness/api"
	"github.com/gorilla/pat"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func parseFlags(args []string) bool {
	flagSet := flag.NewFlagSet("configuration", flag.ExitOnError)
	gVersion := flagSet.Bool("version", false, "Print version and exit")
	flagSet.Parse(args)

	return *gVersion
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get(":path")
	if path == "" {
		path = "index.html"
	}
	fullPath := fmt.Sprintf("dashboard/%s", path)
	data, _ := Asset(fullPath)

	ext := path[strings.LastIndex(path, "."):]
	mimeType := mime.TypeByExtension(ext)
	fmt.Println(path, ext, mimeType)
	w.Header().Set("Content-type", mimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func getRouter() *pat.Router {
	router := pat.New()
	router.Get("/{path:.*}", serveAssets)

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
