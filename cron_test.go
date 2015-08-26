package cron

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cron", func() {
	var subject *cron

	BeforeEach(func() {
		subject = NewCron()
	})

	It("should schedule tasks", func() {
		cnt := 0
		subject.Every(10*time.Millisecond, func() error { cnt++; return nil })
		Eventually(func() int {
			return cnt
		}).Should(BeNumerically(">", 1))
	})

	It("should report errors", func() {
		errs := 0
		subject.OnError(func(err error) { errs++ })
		subject.Every(10*time.Millisecond, func() error { return errors.New("") })
		Eventually(func() int {
			return errs
		}).Should(BeNumerically(">", 1))
	})

	It("should have default cron", func() {
		cnt := 0
		Every(10*time.Millisecond, func() error { cnt++; return nil })
		Eventually(func() int {
			return cnt
		}).Should(BeNumerically(">", 1))
	})

	It("should run task immediately", func() {
		cnt := 0
		subject.Every(time.Hour, func() error { cnt++; return nil })
		Eventually(func() int {
			return cnt
		}).Should(Equal(1))
	})
})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cron")
}
