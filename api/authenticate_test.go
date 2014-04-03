package api

import (
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
	"os"
)

func (s *Suite) TestStrConfigOrEnv(c *gocheck.C) {
	c.Assert(StrConfigOrEnv("key-that-not-exists"), gocheck.Equals, "")

	os.Setenv("my-test-key", "value-of-my-test-key")
	c.Assert(StrConfigOrEnv("my-test-key"), gocheck.Equals, "value-of-my-test-key")

	os.Setenv("my-test-key", "")
	os.Setenv("MY-TEST-KEY", "value-of-my-test-key")
	c.Assert(StrConfigOrEnv("my-test-key"), gocheck.Equals, "value-of-my-test-key")

	os.Setenv("MY-TEST-KEY", "")
	config.Set("my-test-key", "value-of-my-test-key")
	c.Assert(StrConfigOrEnv("my-test-key"), gocheck.Equals, "value-of-my-test-key")
}
