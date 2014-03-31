package main

import (
	"github.com/gorilla/mux"
	"launchpad.net/gocheck"
	"net/http"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {
}

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestLoadConfig(c *gocheck.C) {
	c.Assert(nil, gocheck.IsNil)
}

func routeExists(method string, url string) bool {
	var match *mux.RouteMatch = &mux.RouteMatch{}

	router := getRouter()

	request, _ := http.NewRequest(method, url, nil)

	return router.Match(request, match)
}

func (s *Suite) TestRouterHasHealthcheck(c *gocheck.C) {
	c.Assert(routeExists("GET", "/healthcheck"), gocheck.Equals, true)
}
