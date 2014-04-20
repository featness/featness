package api

import (
	"encoding/json"
	"github.com/globoi/featness/api/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Team Handler", func() {
	var (
		users *mgo.Collection
		conn  *mgo.Session
	)

	BeforeEach(func() {
		var (
			err error
		)

		conn, users, err = models.Users()
		Expect(err).ShouldNot(HaveOccurred())
		users.RemoveAll(bson.M{})
	})

	AfterEach(func() {
		conn.Close()
		conn = nil
	})

	Context("when auto-completing username", func() {

		It("should fail if userid argument not found", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/users/find", nil)
			Expect(err).ShouldNot(HaveOccurred())

			FindUsersWithIdLike(recorder, request)
			Expect(recorder.Code).Should(Equal(http.StatusBadRequest))

			body, err := ioutil.ReadAll(recorder.Body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).Should(BeEmpty())
		})

		It("should return empty array if not found", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/users/find?name=invalid", nil)
			Expect(err).ShouldNot(HaveOccurred())

			FindUsersWithIdLike(recorder, request)
			Expect(recorder.Code).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(recorder.Body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).ShouldNot(BeEmpty())

			var obj []models.User
			err = json.Unmarshal(body, &obj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(obj).To(BeEmpty())
		})

	})
})
