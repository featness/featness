package api

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"launchpad.net/gocheck"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTeamModel(t *testing.T) {
	RegisterFailHandler(Fail)
	MongoStartup("featness-tests", "localhost:3333", "featnesstests", "", "")
	RunSpecs(t, "Team Model Suite")
}

var _ = Describe("Team", func() {
	var (
		conn  *mgo.Database
		teams *mgo.Collection
	)

	BeforeEach(func() {
		conn, _ = Conn()
		teams = conn.C("Team")
		teams.Remove(bson.M{})
	})

	It("Can get all teams for a given member", func() {
		err := teams.Insert(
			&Team{"test1", []string{"heynemann"}},
			&Team{"test2", []string{"john"}},
		)
		Expect(err).ShouldNot(HaveOccurred())

		var team Team
		err = teams.Find(bson.M{"name": "test1"}).One(&team)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(team).ShouldNot(BeNil())
		Expect(team.Name).Should(Equal("test1"))
		Expect(team.Members).Should(HaveLen(1))
		Expect(team.Members[0]).Should(Equal("heynemann"))
	})

})
