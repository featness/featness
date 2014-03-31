package api

import (
	"launchpad.net/gocheck"
)

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestLoadConfig(c *gocheck.C) {
	c.Assert(nil, gocheck.IsNil)
}
