package main

import (
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api"
	"github.com/globoi/featness/api/models"
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

var securityKey string

type SecureFunc func(http.ResponseWriter, *http.Request, *jwt.Token)

func AuthRequiredFunc(handler SecureFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header, ok := r.Header["X-Auth-Token"]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("X-Auth-Token status was not found.")
			return
		}

		token, err := jwt.Parse(header[0], func(t *jwt.Token) ([]byte, error) { return []byte(securityKey), nil })
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Sprintf("X-Auth-Token is not a valid token (%v).", err))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("X-Auth-Token is not a valid token.")
			return
		}

		handler(w, r, token)
	}
}

func getRouter() *pat.Router {
	router := pat.New()
	router.Get("/healthcheck", AllowCrossDomainFunc(api.Healthcheck))
	router.Post("/authenticate/google", AllowCrossDomainFunc(api.AuthenticateWithGoogle))
	router.Post("/authenticate/facebook", AllowCrossDomainFunc(api.AuthenticateWithFacebook))
	router.Get("/authenticate/validate", AllowCrossDomainFunc(api.IsAuthenticationValid([]byte(securityKey))))
	router.Post("/teams/new", AllowCrossDomainFunc(AuthRequiredFunc(api.CreateTeam)))
	router.Get("/teams/available", AllowCrossDomainFunc(AuthRequiredFunc(api.IsTeamNameAvailable)))
	router.Get("/teams", AllowCrossDomainFunc(AuthRequiredFunc(api.GetUserTeams)))
	router.Get("/users/find", AllowCrossDomainFunc(api.FindUsersWithIdLike))
	router.Get("/all-teams", AllowCrossDomainFunc(api.GetAllTeams))
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

	return models.MongoStartup("featness", hosts, database, username, password)
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

	key, err := config.GetString("security_key")
	if err != nil {
		fmt.Printf("There must be an unique security key in the configuration file at %s. None found.\n", configFile)
		return
	}
	securityKey = key

	router := getRouter()
	log.Println("featness-api running at http://localhost:8000...")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Printf("ERROR: Binding to port 8000 failed with error:\n\t>>> %s\n", err)
		return
	}
}
