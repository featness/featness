package main

import (
	"launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {
}

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestLoadConfig(c *gocheck.C) {
	c.Assert("bla", gocheck.IsNil)
}
