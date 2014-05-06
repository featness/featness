package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
	"log"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {
}

var _ = gocheck.Suite(&Suite{})

func loadConfig(path string) {
	err := config.ReadAndWatchConfigFile(path)
	if err != nil {
		msg := `Could not find featness-api config file. Searched on %s.
	For an example conf check api/etc/local.conf file.\n %s`
		log.Panicf(msg, path, err)
	}
}

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	MongoStartup("featness", "localhost:3334", "featness", "", "")
	RunSpecs(t, "API Models Suite")
}
