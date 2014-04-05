package api

import (
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
)

func (s *MongoSuite) TestRouterHasAuthGoogle(c *gocheck.C) {
	teams := s.conn.C("Team")
	defer teams.Remove(bson.M{"name": "test1"})

	err := teams.Insert(&Team{"test1", []string{"heynemann"}})
	c.Assert(err, gocheck.IsNil)

	var team Team
	err = teams.Find(bson.M{"name": "test1"}).One(&team)
	c.Assert(err, gocheck.IsNil)

	c.Assert(team, gocheck.NotNil)
	c.Assert(team.Name, gocheck.Equals, "test1")
	c.Assert(len(team.Members), gocheck.Equals, 1)
	c.Assert(team.Members[0], gocheck.Equals, "heynemann")
}
