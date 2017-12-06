package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfLoupe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CfLoupe Suite")
}
