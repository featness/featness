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
		teams *mgo.Collection
	)

	BeforeEach(func() {
		teams = Teams()
		teams.Remove(bson.M{})
	})

	Context("when the user is in a team", func() {
		It("Can get all teams for a given member", func() {
			err := teams.Insert(
				&Team{"test1", []string{"heynemann"}},
				&Team{"test2", []string{"john"}},
			)
			Expect(err).ShouldNot(HaveOccurred())

			userTeams, err := GetTeamsFor("heynemann")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(userTeams).Should(HaveLen(1))
			Expect(userTeams[0].Name).Should(Equal("test1"))
			Expect(userTeams[0].Members).Should(HaveLen(1))
			Expect(userTeams[0].Members[0]).Should(Equal("heynemann"))
		})
	})

	Context("when the user is in no teams", func() {
		It("should return an empty list of teams", func() {
			err := teams.Insert(
				&Team{"test1", []string{"jane"}},
				&Team{"test2", []string{"john"}},
			)
			Expect(err).ShouldNot(HaveOccurred())

			userTeams, err := GetTeamsFor("heynemann")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(userTeams).Should(BeEmpty())
		})
	})

})
