package api

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
		teams     *mgo.Collection
		users     *mgo.Collection
		testUsers []User
	)

	BeforeEach(func() {
		testUsers = []User{}
		teams = Teams()
		teams.RemoveAll(bson.M{})

		users = Users()
		users.RemoveAll(bson.M{})

		for i := 0; i < 10; i++ {
			userID := "userID" + string(i)
			users.Insert(
				&User{bson.NewObjectId(), "User " + string(i), userID, time.Now(), "http://picture.com/" + string(i)},
			)
			result := User{}
			err := users.Find(bson.M{"userid": userID}).One(&result)
			if err != nil {
				log.Panicf(err.Error())
			}
			testUsers = append(testUsers, result)
		}
	})

	Context(" - User model", func() {
		Context("when the user is logging in for the first time", func() {
			It("Should create a new user and return it", func() {
				user, err := GetOrCreateUser("Bernardo Heynemann", "heynemann", "http://my.picture.url")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Name).Should(Equal("Bernardo Heynemann"))
				Expect(user.UserID).Should(Equal("heynemann"))
				Expect(user.ImageURL).Should(Equal("http://my.picture.url"))
			})
		})

		Context("when the user is logging again", func() {
			It("Shouldn't change the existing ObjectId", func() {
				user := testUsers[0]
				newUser, err := GetOrCreateUser(user.Name, user.UserID, user.ImageURL)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(newUser).ShouldNot(BeNil())
				Expect(user.Id).Should(Equal(newUser.Id))
			})
		})
	})

	Context(" - Team model", func() {
		Context("when the user is in a team", func() {
			It("Can get all teams for a given member", func() {
				err := teams.Insert(
					&Team{bson.NewObjectId(), "test1-team1", []bson.ObjectId{testUsers[0].Id}},
					&Team{bson.NewObjectId(), "test1-team2", []bson.ObjectId{testUsers[1].Id}},
				)
				Expect(err).ShouldNot(HaveOccurred())

				userTeams, err := GetTeamsFor(testUsers[0].Id)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(userTeams).Should(HaveLen(1))
				Expect(userTeams[0].Name).Should(Equal("test1-team1"))
				Expect(userTeams[0].Members).Should(HaveLen(1))
				Expect(userTeams[0].Members[0]).Should(Equal(testUsers[0].Id))
			})
		})

		Context("when the user is in no teams", func() {
			It("should return an empty list of teams", func() {
				err := teams.Insert(
					&Team{bson.NewObjectId(), "test2-team1", []bson.ObjectId{testUsers[2].Id}},
					&Team{bson.NewObjectId(), "test2-team2", []bson.ObjectId{testUsers[3].Id}},
				)
				Expect(err).ShouldNot(HaveOccurred())

				userTeams, err := GetTeamsFor(testUsers[4].Id)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(userTeams).Should(BeEmpty())
			})
		})
	})

})
