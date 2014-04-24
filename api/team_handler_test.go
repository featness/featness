package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Team Handler", func() {
	var (
		teams *mgo.Collection
		users *mgo.Collection
		user  *models.User
		conn  *mgo.Session
	)

	BeforeEach(func() {
		var (
			err error
		)

		conn, teams, err = models.Teams()
		Expect(err).ShouldNot(HaveOccurred())
		teams.RemoveAll(bson.M{})

		conn, users, err = models.UsersWithConn(conn)
		Expect(err).ShouldNot(HaveOccurred())
		users.RemoveAll(bson.M{})
	})

	AfterEach(func() {
		conn.Close()
		conn = nil
	})

	Context("when no teams registered", func() {

		Context("when obtaining user teams", func() {

			It("should fail if user not found", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = "invalid-user-id"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				GetUserTeams(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusBadRequest))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).Should(BeEmpty())
			})

			It("should return user teams with an empty array", func() {
				userID := "user-1"
				user = &models.User{bson.NewObjectId(), "facebook", "token", "User 1", userID, time.Now(), "http://picture.com/1"}
				users.Insert(user)

				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = userID
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				GetUserTeams(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).ShouldNot(BeNil())

				var obj []interface{}
				err = json.Unmarshal(body, &obj)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(obj).To(HaveLen(0))
			})
		})

		Context("when obtaining all teams", func() {
			It("should return all teams with an empty array", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/all-teams", nil)
				Expect(err).ShouldNot(HaveOccurred())

				GetAllTeams(recorder, request)
				Expect(recorder.Code).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).ShouldNot(BeNil())

				var obj []interface{}
				err = json.Unmarshal(body, &obj)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(obj).To(HaveLen(0))
			})
		})

		Context("when verifying if a name is available", func() {
			It("should return true for any names", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams/available?name=something", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = "user-1"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				IsTeamNameAvailable(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).ShouldNot(BeNil())

				var obj interface{}
				err = json.Unmarshal(body, &obj)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(obj.(bool)).To(BeTrue())
			})

			It("Should fail when no name provided", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams/available", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = "user-1"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				IsTeamNameAvailable(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("when teams registered", func() {
		var (
			allUsers []*models.User
			allTeams []*models.Team
		)

		BeforeEach(func() {
			allUsers = []*models.User{}
			allTeams = []*models.Team{}

			for i := 0; i < 10; i++ {
				userID := "userID" + string(i)
				users.Insert(
					&models.User{bson.NewObjectId(), "facebook", "token", "User " + string(i), userID, time.Now(), "http://picture.com/" + string(i)},
				)
				result := models.User{}
				err := users.Find(bson.M{"userid": userID}).One(&result)
				if err != nil {
					log.Panicf(err.Error())
				}
				allUsers = append(allUsers, &result)
			}

			team, _ := models.GetOrCreateTeam("team1", *allUsers[0], *allUsers[1])
			team2, _ := models.GetOrCreateTeam("team2", *allUsers[2], *allUsers[3])
			allTeams = append(allTeams, team, team2)
		})

		It("Should get all teams", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/all-teams", nil)
			Expect(err).ShouldNot(HaveOccurred())

			GetAllTeams(recorder, request)
			Expect(recorder.Code).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(recorder.Body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).ShouldNot(BeNil())

			var obj []interface{}
			err = json.Unmarshal(body, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(obj).To(HaveLen(2))
		})

		It("should return user teams when available", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/teams", nil)
			Expect(err).ShouldNot(HaveOccurred())

			token := jwt.New(jwt.GetSigningMethod("HS256"))
			token.Claims["sub"] = allUsers[1].UserId
			token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			GetUserTeams(recorder, request, token)
			Expect(recorder.Code).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(recorder.Body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).ShouldNot(BeNil())

			var obj []models.Team
			err = json.Unmarshal(body, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(obj).To(HaveLen(1))
			Expect(obj[0].Name).To(Equal("team1"))
		})

		Context("when verifying if a name is available", func() {
			It("should return false for taken names", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams/available?name=team1", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = "user-1"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				IsTeamNameAvailable(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).ShouldNot(BeNil())

				var obj interface{}
				err = json.Unmarshal(body, &obj)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(obj.(bool)).To(BeFalse())
			})

			It("should return true for available names", func() {
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest("GET", "/teams/available?name=team3", nil)
				Expect(err).ShouldNot(HaveOccurred())

				token := jwt.New(jwt.GetSigningMethod("HS256"))
				token.Claims["sub"] = "user-1"
				token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				IsTeamNameAvailable(recorder, request, token)
				Expect(recorder.Code).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(body)).ShouldNot(BeNil())

				var obj interface{}
				err = json.Unmarshal(body, &obj)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(obj.(bool)).To(BeTrue())
			})
		})
	})
})
