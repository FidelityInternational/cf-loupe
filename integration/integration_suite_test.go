package integration_test

import (
	"fmt"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var session *gexec.Session

var _ = BeforeSuite(func() {
	checkVar("CF_FOUNDATION_1")
	checkVar("CF_API_1")
	checkVar("CF_USERNAME_1")
	checkVar("CF_PASSWORD_1")
	pathToCFLoupe, err := gexec.Build("github.com/FidelityInternational/cf-loupe")
	Expect(err).To(Succeed())

	cmd := exec.Command(pathToCFLoupe)
	cmd.Dir = "../"
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).To(Succeed())

	Eventually(session, 10).Should(gbytes.Say("Starting app on port"))
})

var _ = AfterSuite(func() {
	if session != nil {
		session.Kill()
	}
})

func checkVar(varName string) {
	defer GinkgoRecover()
	if os.Getenv(varName) == "" {
		Fail(fmt.Sprintf("%s is required but was not set", varName))
	}
}
