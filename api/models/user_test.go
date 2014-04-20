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

		conn, teams, err = Teams()
		Expect(err).ShouldNot(HaveOccurred())
		teams.RemoveAll(bson.M{})

		testUsers = []User{}
		conn, users, err = UsersWithConn(conn)
		Expect(err).ShouldNot(HaveOccurred())

		users.RemoveAll(bson.M{})

		for i := 0; i < 10; i++ {
			userID := "userID" + string(i)
			users.Insert(
				&User{bson.NewObjectId(), "facebook", "token", userID, "User " + string(i), time.Now(), "http://picture.com/" + string(i)},
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

	Context(" - User model", func() {
		Context("when the user is logging in for the first time", func() {
			It("Should create a new user and return it", func() {
				user, err := GetOrCreateUser("facebook", "token", "heynemann", "Bernardo Heynemann", "http://my.picture.url")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Name).Should(Equal("Bernardo Heynemann"))
				Expect(user.AccessToken).Should(Equal("token"))
				Expect(user.UserId).Should(Equal("heynemann"))
				Expect(user.ImageUrl).Should(Equal("http://my.picture.url"))
			})
		})

		Context("when the user is logging again", func() {
			It("Shouldn't change the existing ObjectId", func() {
				user := testUsers[0]
				newUser, err := GetOrCreateUser(user.Provider, user.AccessToken, user.UserId, user.Name, user.ImageUrl)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(newUser).ShouldNot(BeNil())
				Expect(user.Id).Should(Equal(newUser.Id))
			})
		})

		Context("when autocompleting user by name", func() {

			It("should return an empty array if not users found", func() {
				users, err := FindUsersWithIdLike("invalid id")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(users).Should(BeEmpty())
			})

		})
	})

})
