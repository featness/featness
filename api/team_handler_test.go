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
		teams *mgo.Collection
	)

	BeforeEach(func() {
		teams = models.Teams()
		teams.RemoveAll(bson.M{})
	})

	Context("when no teams registered", func() {

		It("should return user teams", func() {
		})

		It("should return all teams", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/all-teams", nil)
			Expect(err).ShouldNot(HaveOccurred())

			GetAllTeamsUnauthenticated(recorder, request)
			Expect(recorder.Code).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(recorder.Body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).ShouldNot(BeNil())

			obj, err := json.Marshal(body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(obj).To(HaveLen(0))
		})
	})

	Context("when teams registered", func() {

	})

})
