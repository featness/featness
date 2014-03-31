package api

import (
	"io/ioutil"
	"launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestHealthcheck(c *gocheck.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/healthcheck", nil)
	c.Assert(err, gocheck.IsNil)

	Healthcheck(recorder, request)
	c.Assert(recorder.Code, gocheck.Equals, http.StatusOK)

	body, err := ioutil.ReadAll(recorder.Body)
	c.Assert(err, gocheck.IsNil)
	c.Assert(string(body), gocheck.Equals, "WORKING")
}
