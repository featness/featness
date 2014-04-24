package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

var _ = Describe("Models", func() {
	var (
		conn      *mgo.Session
		teams     *mgo.Collection
		users     *mgo.Collection
		testUsers []User
	)

	BeforeEach(func() {
		var (
			err error
		)

		testUsers = []User{}

		conn, teams, err = Teams()
		Expect(err).ShouldNot(HaveOccurred())

		teams.RemoveAll(bson.M{})

		conn, users, err = UsersWithConn(conn)
		Expect(err).ShouldNot(HaveOccurred())

		users.RemoveAll(bson.M{})

		for i := 0; i < 10; i++ {
			userID := "userID" + string(i)
			users.Insert(
				&User{bson.NewObjectId(), "facebook", "token", "User " + string(i), userID, time.Now(), "http://picture.com/" + string(i)},
			)
			result := User{}
			err := users.Find(bson.M{"userid": userID}).One(&result)
			if err != nil {
				log.Panicf(err.Error())
			}
			testUsers = append(testUsers, result)
		}
	})

	AfterEach(func() {
		conn.Close()
		conn = nil
	})

	Context(" - Team model", func() {
		Context("when the user is in a team", func() {
			It("Can get all teams for a given member", func() {
				_, err := GetOrCreateTeam("test1-team1", testUsers[0], testUsers[1])
				_, err2 := GetOrCreateTeam("test1-team2", testUsers[2], testUsers[3], testUsers[4])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(err2).ShouldNot(HaveOccurred())

				userTeams, err := GetTeamsFor(testUsers[0].Id)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(userTeams).Should(HaveLen(1))
				Expect(userTeams[0].Name).Should(Equal("test1-team1"))
				Expect(userTeams[0].Owner).Should(Equal(testUsers[0].Id))
				Expect(userTeams[0].Members).Should(HaveLen(1))
				Expect(userTeams[0].Members[0]).Should(Equal(testUsers[0].Id))
			})

			It("Can get all teams for a given owner", func() {
				_, err := GetOrCreateTeam("test4-team1", testUsers[0], testUsers[1])
				_, err2 := GetOrCreateTeam("test4-team2", testUsers[2], testUsers[3], testUsers[4])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(err2).ShouldNot(HaveOccurred())

				userTeams, err := GetTeamsFor(testUsers[0].Id)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(userTeams).Should(HaveLen(1))
				Expect(userTeams[0].Name).Should(Equal("test4-team1"))
				Expect(userTeams[0].Owner).Should(Equal(testUsers[0].Id))
				Expect(userTeams[0].Members).Should(HaveLen(1))
				Expect(userTeams[0].Members[0]).Should(Equal(testUsers[0].Id))
			})

		})

		Context("when the user is in no teams", func() {
			It("should return an empty list of teams", func() {
				_, err := GetOrCreateTeam("test2-team1", testUsers[2])
				_, err2 := GetOrCreateTeam("test2-team2", testUsers[3], testUsers[4])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(err2).ShouldNot(HaveOccurred())

				userTeams, err := GetTeamsFor(testUsers[5].Id)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(userTeams).Should(BeEmpty())
			})
		})

		Context("when the team doesn't exist", func() {
			It("should create team", func() {
				team, err := GetOrCreateTeam("team1", testUsers[6])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(team.Name).Should(Equal("team1"))
				Expect(team.Owner).Should(Equal(testUsers[6].Id))
				Expect(team.Members).Should(HaveLen(0))
			})

			It("should create team with user", func() {
				team, err := GetOrCreateTeam("team1", testUsers[0], testUsers[1])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(team.Name).Should(Equal("team1"))
				Expect(team.Owner).Should(Equal(testUsers[0].Id))
				Expect(team.Members).Should(HaveLen(1))
				Expect(team.Members[0]).Should(Equal(testUsers[1].Id))
			})

			It("should create team with users", func() {
				team, err := GetOrCreateTeam("team1", testUsers[0], testUsers[1], testUsers[2])

				Expect(err).ShouldNot(HaveOccurred())
				Expect(team.Name).Should(Equal("team1"))
				Expect(team.Members).Should(HaveLen(2))
			})
		})
	})

})
