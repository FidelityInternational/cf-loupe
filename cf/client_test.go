package cf_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/FidelityInternational/cf-loupe/helpers"

	. "github.com/FidelityInternational/cf-loupe/cf"
)

var _ = Describe("BuildClientConfigFromEnvironment", func() {
	var fakeEnv []string
	var fapi1 *helpers.FakeApi
	var fapi2 *helpers.FakeApi

	BeforeEach(func() {
		fapi1 = helpers.NewFakeApi()
		fapi2 = helpers.NewFakeApi()
	})

	AfterEach(func() {
		fapi1.TeardownFakeApi()
		fapi2.TeardownFakeApi()
	})

	Context("When a foundation is missing some credentials", func() {
		BeforeEach(func() {
			fakeEnv = []string{
				"CF_USERNAME_1=admin",
				"CF_PASSWORD_1=1234",
				"CF_FOUNDATION_1=dev",
				fmt.Sprintf("CF_API_1=%s", fapi1.Server.URL),
			}
		})

		It("returns the correct map of client configurations", func() {
			foundationConfigs, err := BuildClientConfigFromEnvironment(fakeEnv)
			Expect(err).To(Succeed())
			Expect(foundationConfigs).To(HaveKey("dev"))

			devConfig := foundationConfigs["dev"]
			Expect(devConfig.Username).To(Equal("admin"))
			Expect(devConfig.Password).To(Equal("1234"))
			Expect(devConfig.ApiAddress).To(Equal(fapi1.Server.URL))
		})

		It("returns the correct map of clients", func() {
			clients, err := BuildClientsFromEnvironment(fakeEnv)
			Expect(err).To(Succeed())
			Expect(clients).To(HaveKey("dev"))
		})
	})

	Context("When there are no foundations", func() {
		BeforeEach(func() {
			fakeEnv = []string{}
		})

		It("returns a meaningful error", func() {
			_, err := BuildClientConfigFromEnvironment(fakeEnv)
			Expect(err).To(MatchError("no foundation environment variables found. CF_USERNAME_1, CF_PASSWORD_1, CF_API_1 and CF_FOUNDATION_1 must be set"))
		})
	})

	Context("When a foundation is missing some credentials", func() {
		BeforeEach(func() {
			fakeEnv = []string{
				"CF_USERNAME_1=admin",
				"CF_PASSWORD_1=1234",
				"CF_FOUNDATION_1=dev",
			}
		})

		It("returns a meaningful error", func() {
			_, err := BuildClientConfigFromEnvironment(fakeEnv)
			Expect(err).To(MatchError("CF_API_1 env var not found for dev foundation"))
		})
	})

	Context("When there are multiple foundations", func() {
		BeforeEach(func() {
			fakeEnv = []string{
				"SHELL=/bin/zsh",
				"TERM=xterm-256color",
				"CF_USERNAME_1=admin",
				"CF_PASSWORD_1=1234",
				"CF_FOUNDATION_1=dev",
				fmt.Sprintf("CF_API_1=%s", fapi1.Server.URL),
				"CF_USERNAME_2=admin",
				"CF_PASSWORD_2=12345",
				"CF_FOUNDATION_2=prod",
				fmt.Sprintf("CF_API_2=%s", fapi2.Server.URL),
			}
		})

		It("returns the correct map of client configurations", func() {
			foundationConfigs, err := BuildClientConfigFromEnvironment(fakeEnv)
			Expect(err).To(Succeed())

			Expect(foundationConfigs).To(HaveKey("dev"))
			devConfig := foundationConfigs["dev"]
			Expect(devConfig.Username).To(Equal("admin"))
			Expect(devConfig.Password).To(Equal("1234"))
			Expect(devConfig.ApiAddress).To(Equal(fapi1.Server.URL))

			Expect(foundationConfigs).To(HaveKey("prod"))
			prodConfig := foundationConfigs["prod"]
			Expect(prodConfig.Username).To(Equal("admin"))
			Expect(prodConfig.Password).To(Equal("12345"))
			Expect(prodConfig.ApiAddress).To(Equal(fapi2.Server.URL))
		})

		It("returns the correct map of clients", func() {
			clients, err := BuildClientsFromEnvironment(fakeEnv)
			Expect(err).To(Succeed())
			Expect(clients).To(HaveKey("dev"))
			Expect(clients).To(HaveKey("prod"))
		})
	})
})
