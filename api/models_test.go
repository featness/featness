package api

import (
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
)

func (s *MongoSuite) TestRouterHasAuthGoogle(c *gocheck.C) {
	teams := s.conn.C("Team")
	defer teams.Remove(bson.M{"Name": "test1"})
}
