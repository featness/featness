package main

import (
	"github.com/globoi/featness/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"github.com/maraino/go-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tsuru/config"
	"net/http"
	"testing"
)

var (
	router *pat.Router
)

var _ = Describe("API Server Main Module", func() {

	BeforeEach(func() {
		router = getRouter()
	})

	Context("when configuring routes", func() {

		It("should have /healthcheck route", func() {
			exists := routeExists("GET", "/healthcheck")
			Expect(exists).Should(BeTrue())
		})

		It("should have /authenticate/google route", func() {
			exists := routeExists("POST", "/authenticate/google")
			Expect(exists).Should(BeTrue())
		})

		It("should have /authenticate/facebook route", func() {
			exists := routeExists("POST", "/authenticate/facebook")
			Expect(exists).Should(BeTrue())
		})

	})

	Context("when loading configuration", func() {
		It("should load configuration keys in the given file", func() {
			logger := &LoggerTest{}
			logger.When("Panicf").Times(0)

			loadConfigFile("../testdata/etc/featness-api1.conf", logger)

			value, errorGetBool := config.GetBool("my_data")
			ok, errorMock := logger.Verify()

			Expect(errorGetBool).ShouldNot(HaveOccurred())
			Expect(errorMock).ShouldNot(HaveOccurred())
			Expect(value).Should(BeTrue())
			Expect(ok).Should(BeTrue())
		})

		It("should fail when given wrong path", func() {
			logger := &LoggerTest{}
			logger.When("Panicf", mock.Any).Times(1)

			loadConfigFile("wrong-path.conf", logger)

			ok, errorMock := logger.Verify()
			Expect(ok).Should(BeTrue())
			Expect(errorMock).Should(BeNil())
		})
	})

	Context("when parsing flags", func() {
		It("should get configuration file path", func() {
			configFile, gVersion := parseFlags([]string{"--config", "my.conf"})

			Expect(configFile).Should(Equal("my.conf"))
			Expect(gVersion).Should(BeFalse())
		})

		It("should get app version", func() {
			configFile, gVersion := parseFlags([]string{"--version"})

			Expect(configFile).Should(Equal("/etc/featness-api.conf"))
			Expect(gVersion).Should(BeTrue())
		})

	})

})

func routeExists(method string, url string) bool {
	var match *mux.RouteMatch = &mux.RouteMatch{}

	request, _ := http.NewRequest(method, url, nil)

	return router.Match(request, match)
}

type LoggerTest struct {
	mock.Mock
}

func (l *LoggerTest) Panicf(format string, v ...interface{}) {
	l.Called(format)
}

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	api.MongoStartup("featness", "localhost:3334", "featness", "", "")
	RunSpecs(t, "API Server Main Suite")
}
