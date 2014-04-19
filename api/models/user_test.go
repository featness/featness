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
				&User{bson.NewObjectId(), "facebook", userID, "User " + string(i), time.Now(), "http://picture.com/" + string(i)},
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
				user, err := GetOrCreateUser("facebook", "heynemann", "Bernardo Heynemann", "http://my.picture.url")
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
				newUser, err := GetOrCreateUser(user.Provider, user.UserID, user.Name, user.ImageURL)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(newUser).ShouldNot(BeNil())
				Expect(user.Id).Should(Equal(newUser.Id))
			})
		})
	})

})
