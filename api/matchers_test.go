package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Matcher Methods", func() {

	var (
		matcherPaths []string
	)

	Context("when no matcher is registered", func() {
		BeforeEach(func() {
			matcherPaths = []string{}
		})

		It("should return false", func() {
			data := &RequestData{}
			Expect(Matches(matcherPaths, data)).Should(BeTrue())
		})
	})

	Context("when a single matcher is registered", func() {
		BeforeEach(func() {
			matcherPaths = []string{
				"../testdata/matcher.sh",
			}
		})

		It("should return true", func() {
			data := &RequestData{}
			Expect(Matches(matcherPaths, data)).Should(BeTrue())
		})
	})

	Context("when multiple true matchers are registered", func() {
		BeforeEach(func() {
			matcherPaths = []string{
				"../testdata/matcher.sh",
				"../testdata/matcher.sh",
			}
		})

		It("should return true", func() {
			data := &RequestData{}
			Expect(Matches(matcherPaths, data)).Should(BeTrue())
		})
	})

	Context("when multiple mixed value matchers are registered", func() {
		BeforeEach(func() {
			matcherPaths = []string{
				"../testdata/matcher.sh",
				"../testdata/matcher-false.sh",
				"../testdata/matcher.sh",
				"../testdata/matcher-false.sh",
				"../testdata/matcher.sh",
				"../testdata/matcher-very-slow.sh",
			}
		})

		It("should return true", func() {
			data := &RequestData{}
			Expect(Matches(matcherPaths, data)).Should(BeTrue())
		})
	})

	Context("when only false value matchers are registered", func() {
		BeforeEach(func() {
			matcherPaths = []string{
				"../testdata/matcher-false.sh",
				"../testdata/matcher-false.sh",
			}
		})

		It("should return false", func() {
			data := &RequestData{}
			Expect(Matches(matcherPaths, data)).Should(BeFalse())
		})
	})
})
