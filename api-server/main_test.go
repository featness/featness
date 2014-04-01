package main

import (
	"github.com/gorilla/mux"
	"github.com/maraino/go-mock"
	"launchpad.net/gocheck"
	"net/http"
	"github.com/tsuru/config"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {}

var _ = gocheck.Suite(&Suite{})

type LoggerTest struct {
    mock.Mock
}

func (l *LoggerTest) Panicf(format string, v ...interface{}) {
    l.Called(format, )
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

func (s *Suite) TestLoadConfig(c *gocheck.C) {
	logger := &LoggerTest{}
	logger.When("Panicf").Times(0)

	loadConfigFile("../testdata/etc/featness-api1.conf", logger)
	
	value, errorGetBool := config.GetBool("my_data")
	ok, errorMock := logger.Verify()
	
	c.Assert(value, gocheck.Equals, true)
	c.Assert(errorGetBool, gocheck.IsNil)
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(errorMock, gocheck.IsNil)
}

func (s *Suite) TestLoadConfigWhenWrongPath(c *gocheck.C) {
	logger := &LoggerTest{}
	logger.When("Panicf", `Could not find featness-api config file. Searched on %s.
	For an example conf check featness-api/etc/featness-api.conf file.\n %s`).Times(1)

	loadConfigFile("wrong-path.conf", logger)

	ok, errorMock := logger.Verify()
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(errorMock, gocheck.IsNil)
}