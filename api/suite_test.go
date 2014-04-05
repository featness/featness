package api

import (
	"github.com/tsuru/config"
	"labix.org/v2/mgo"
	"launchpad.net/gocheck"
	"log"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {
}

var _ = gocheck.Suite(&Suite{})

func loadConfig(path string) {
	err := config.ReadAndWatchConfigFile(path)
	if err != nil {
		msg := `Could not find featness-api config file. Searched on %s.
	For an example conf check api/etc/local.conf file.\n %s`
		log.Panicf(msg, path, err)
	}
}

type MongoSuite struct {
	session *mgo.Session
	conn    *mgo.Database
}

var _ = gocheck.Suite(&MongoSuite{})

func (s *MongoSuite) SetUpSuite(c *gocheck.C) {
	MongoStartup("featness-tests", "localhost:3333", "featnesstests", "", "")
}

func (s *MongoSuite) SetUpTest(c *gocheck.C) {
	session, err := CopyMonotonicSession("featness-tests")
	c.Assert(err, gocheck.IsNil)

	s.session = session
	s.conn = session.DB("featnesstests")
}

func (s *MongoSuite) TearDownTest(c *gocheck.C) {
	s.session.Close()
}
