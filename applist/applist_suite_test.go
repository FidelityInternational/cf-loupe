package applist_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApplist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Applist Suite")
}
