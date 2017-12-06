package main_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	xmlpath "gopkg.in/xmlpath.v2"

	gocf "github.com/cloudfoundry-community/go-cfclient"

	. "github.com/FidelityInternational/cf-loupe"
	"github.com/FidelityInternational/cf-loupe/applist"
	"github.com/FidelityInternational/cf-loupe/cf"
	"github.com/FidelityInternational/cf-loupe/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeClient struct {
	ReAuthFunc        func() error
	ListAppsFunc      func() ([]gocf.App, error)
	GetBuildpacksFunc func() (map[string]gocf.Buildpack, error)
	GetOrgsFunc       func() (map[string]gocf.Org, error)
	GetSpacesFunc     func() (map[string]gocf.Space, error)
}

func (client FakeClient) ReAuth() error {
	return client.ReAuthFunc()
}

func (client FakeClient) ListApps() ([]gocf.App, error) {
	return client.ListAppsFunc()
}

func (client FakeClient) GetBuildpacks() (map[string]gocf.Buildpack, error) {
	return client.GetBuildpacksFunc()
}

func (client FakeClient) GetOrgs() (map[string]gocf.Org, error) {
	return client.GetOrgsFunc()
}

func (client FakeClient) GetSpaces() (map[string]gocf.Space, error) {
	return client.GetSpacesFunc()
}

var _ = Describe("Main", func() {
	var server *httptest.Server
	var cfClient FakeClient
	var realCfClient cf.IClient
	var fakeEnv []string
	var fakeApi *helpers.FakeApi

	BeforeEach(func() {
		timeNow := func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2017-08-15T15:00:06Z")
			return t
		}

		cfClients := map[string]cf.IClient{
			"dev": &cfClient,
		}

		server = httptest.NewServer(BuildRouter(cfClients, timeNow))

		fakeApi = helpers.NewFakeApi()
		fakeEnv = []string{
			"CF_USERNAME_1=admin",
			"CF_PASSWORD_1=1234",
			"CF_FOUNDATION_1=dev",
			fmt.Sprintf("CF_API_1=%s", fakeApi.Server.URL),
		}
		cfRealClients, err := cf.BuildClientsFromEnvironment(fakeEnv)
		Expect(err).To(Succeed())
		realCfClient = cfRealClients["dev"]

		cfClient.ReAuthFunc = func() error {
			return realCfClient.ReAuth()
		}

		cfClient.ListAppsFunc = func() ([]gocf.App, error) {
			return []gocf.App{
				gocf.App{
					Name:                  "app1",
					UpdatedAt:             "2017-08-12T16:41:45Z",
					DetectedBuildpackGuid: "abc123",
					SpaceGuid:             "aaaaa",
					Instances:             1,
					Memory:                64,
					State:                 "started",
				},
				gocf.App{
					Name:                  "app2",
					UpdatedAt:             "2016-07-19T16:41:45Z",
					DetectedBuildpackGuid: "def456",
					SpaceGuid:             "aaaaa",
					Instances:             2,
					Memory:                512,
					State:                 "stopped",
				},
				gocf.App{
					Name:                  "app3",
					UpdatedAt:             "2016-07-28T16:41:45Z",
					DetectedBuildpackGuid: "",
					Buildpack:             "https://github.com/cloudfoundry/staticfile-buildpack",
					SpaceGuid:             "bbbbb",
					Instances:             3,
					Memory:                2048,
					State:                 "started",
				},
			}, nil
		}

		cfClient.GetBuildpacksFunc = func() (map[string]gocf.Buildpack, error) {
			return map[string]gocf.Buildpack{
				"abc123": gocf.Buildpack{
					Name:     "ruby_buildpack",
					Filename: "ruby_buildpack-cached-v1.6.47.zip",
				},
				"def456": gocf.Buildpack{
					Name:     "java_buildpack",
					Filename: "java-buildpack-v1_19-fidelity-abc1234.zip",
				},
				"hij789": gocf.Buildpack{
					Name:     "ruby_buildpack",
					Filename: "ruby_buildpack-cached-v2.0.0.zip",
				},
				"33333": gocf.Buildpack{
					Name:     "ruby_buildpack",
					Filename: "ruby_buildpack-cached-v2.0.1.zip",
				},
			}, nil
		}

		cfClient.GetOrgsFunc = func() (map[string]gocf.Org, error) {
			return map[string]gocf.Org{
				"123123123": gocf.Org{
					Name: "project-x",
				},
			}, nil
		}

		cfClient.GetSpacesFunc = func() (map[string]gocf.Space, error) {
			return map[string]gocf.Space{
				"aaaaa": gocf.Space{
					Name:             "dev",
					OrganizationGuid: "123123123",
				},
				"bbbbb": gocf.Space{
					Name:             "test",
					OrganizationGuid: "123123123",
				},
			}, nil
		}
	})

	Describe("GET /", func() {
		Context("When the API returns an error", func() {
			BeforeEach(func() {
				cfClient.ListAppsFunc = func() ([]gocf.App, error) {
					return nil, errors.New("The server is on fire!")
				}
			})

			It("returns 200 as it is a static page", func() {
				resp, err := http.Get(server.URL)
				Expect(err).To(Succeed())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		It("returns 200 OK", func() {
			resp, err := http.Get(server.URL)
			Expect(err).To(Succeed())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("contains cf-loupe in the header", func() {
			resp, err := http.Get(server.URL)
			Expect(err).To(Succeed())

			doc, err := xmlpath.ParseHTML(resp.Body)
			Expect(err).To(Succeed())

			h1Path := xmlpath.MustCompile("/html/body/section/div/div/div[1]/div[2]/h1")
			header, _ := h1Path.String(doc)
			Expect(header).To(ContainSubstring("CF Loupe"))
		})

		PIt("shows some relevant stats", func() {
			resp, err := http.Get(server.URL)
			Expect(err).To(Succeed())

			doc, err := xmlpath.ParseHTML(resp.Body)
			Expect(err).To(Succeed())

			totalAppsPath := xmlpath.MustCompile("/html/body/div/nav/div[1]/div/p[2]")
			totalApps, _ := totalAppsPath.String(doc)

			Expect(totalApps).To(ContainSubstring("3"))

			staleAppsPath := xmlpath.MustCompile("/html/body/div/nav/div[2]/div/p[2]")
			staleApps, _ := staleAppsPath.String(doc)

			Expect(staleApps).To(ContainSubstring("2"))

			depreactedAppsPath := xmlpath.MustCompile("/html/body/div/nav/div[3]/div/p[2]")
			depreactedApps, _ := depreactedAppsPath.String(doc)

			Expect(depreactedApps).To(ContainSubstring("2"))
		})

		PIt("contains a list of deployed apps and last updated times", func() {
			resp, err := http.Get(server.URL)
			Expect(err).To(Succeed())

			doc, err := xmlpath.ParseHTML(resp.Body)
			Expect(err).To(Succeed())

			xpath := func(path string) string {
				val, _ := xmlpath.MustCompile(path).String(doc)
				return val
			}

			// app names
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[1]")).To(Equal("app1"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[1]")).To(Equal("app2"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[1]")).To(Equal("app3"))

			// foundation name
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[2]")).To(Equal("dev"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[2]")).To(Equal("dev"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[2]")).To(Equal("dev"))

			// org names
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[3]")).To(Equal("project-x"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[3]")).To(Equal("project-x"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[3]")).To(Equal("project-x"))

			// space names
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[4]")).To(Equal("dev"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[4]")).To(Equal("dev"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[4]")).To(Equal("test"))

			// instances
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[5]")).To(Equal("1"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[5]")).To(Equal("2"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[5]")).To(Equal("3"))

			// memory
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[6]")).To(Equal("64"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[6]")).To(Equal("512"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[6]")).To(Equal("2048"))

			// state
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[7]")).To(Equal("started"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[7]")).To(Equal("stopped"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[7]")).To(Equal("started"))

			// updated at
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[8]")).To(Equal("2017-08-12"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[8]")).To(Equal("2016-07-19"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[8]")).To(Equal("2016-07-28"))

			// up-to-date
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[9]")).To(Equal("yes"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[9]")).To(Equal("no"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[9]")).To(Equal("no"))

			// buildpack name and version
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[10]")).To(ContainSubstring("ruby 1.6.47"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[10]")).To(ContainSubstring("java 1.19"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[10]")).To(ContainSubstring("https://github.com/cloudfoundry/staticfile-buildpack"))

			// support status
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[11]")).To(Equal("no"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[11]")).To(Equal("yes"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[11]")).To(Equal("no"))

			// status
			Expect(xpath("/html/body/div/table/tbody/tr[1]/td[12]")).To(Equal("✘"))
			Expect(xpath("/html/body/div/table/tbody/tr[2]/td[12]")).To(Equal("✘"))
			Expect(xpath("/html/body/div/table/tbody/tr[3]/td[12]")).To(Equal("✘"))
		})
	})

	Describe("GET /listapps", func() {
		var url *url.URL

		BeforeEach(func() {
			var err error
			url, err = url.Parse(server.URL)
			Expect(err).To(Succeed())
			url.Path = "/listapps"
		})

		Context("When the endpoint returns an error", func() {
			BeforeEach(func() {
				cfClient.ListAppsFunc = func() ([]gocf.App, error) {
					return nil, errors.New("The server is on fire!")
				}
			})

			It("returns 500 Internal Server Error", func() {
				resp, err := http.Get(url.String())
				Expect(err).To(Succeed())

				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})

			It("returns an error message", func() {
				resp, err := http.Get(url.String())
				Expect(err).To(Succeed())
				defer resp.Body.Close()

				bytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(Succeed())
				defer resp.Body.Close()

				str := string(bytes)
				Expect(str).To(ContainSubstring("The server is on fire!"))
			})
		})

		Context("When the endpoint works", func() {
			It("returns 200 OK", func() {
				resp, err := http.Get(url.String())
				Expect(err).To(Succeed())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns a valid json", func() {
				resp, err := http.Get(url.String())
				Expect(err).To(Succeed())

				bytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(Succeed())
				defer resp.Body.Close()

				var jsonData interface{}
				err = json.Unmarshal(bytes, &jsonData)
				Expect(err).To(Succeed())
			})

			It("response successfully unmarshals to an AppData object", func() {
				resp, err := http.Get(url.String())
				Expect(err).To(Succeed())

				bytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(Succeed())
				defer resp.Body.Close()

				var appData applist.AppData
				err = json.Unmarshal(bytes, &appData)
				Expect(err).To(Succeed())

				Expect(appData.Apps[0].Name).To(Equal("app1"))
				Expect(appData.Apps[0].UpdatedAt).To(Equal("2017-08-12"))
				Expect(appData.Apps[0].IsStale).To(BeFalse())
				Expect(appData.Apps[0].Buildpack.Name).To(Equal("ruby"))
				Expect(appData.Apps[0].Buildpack.Version).To(Equal("1.6.47"))
				Expect(appData.Apps[0].Buildpack.Freshness).To(Equal(2))
				Expect(appData.Apps[0].Buildpack.IsDeprecated).To(BeTrue())
				Expect(appData.Apps[1].Name).To(Equal("app2"))
				Expect(appData.Apps[1].UpdatedAt).To(Equal("2016-07-19"))
				Expect(appData.Apps[1].IsStale).To(BeTrue())
				Expect(appData.Apps[1].Buildpack.Name).To(Equal("java"))
				Expect(appData.Apps[1].Buildpack.Version).To(Equal("1.19"))
				Expect(appData.Apps[1].Buildpack.Freshness).To(Equal(0))
				Expect(appData.Apps[1].Buildpack.IsDeprecated).To(BeFalse())
				Expect(appData.Apps).To(HaveLen(3))
			})
		})

		Context("When the auth token expires", func() {

			BeforeEach(func() {
				// Let the token expire
				time.Sleep(1 * time.Second)
			})

			It("Refreshes token before attempting to listapps", func() {
				fakeApi.MaxTokenRefresh = 1
				fakeApi.TokenExpiresIn = 12 // Set it to > 10s or Token.Valid() will say false

				previousTokenCounter := fakeApi.TokenCounter

				_, err := http.Get(url.String())
				Expect(err).To(Succeed())
				Expect(fakeApi.TokenCounter).To(Equal(previousTokenCounter + 1))
				Expect(fakeApi.TokenRefreshCounter).To(Equal(1))
			})

			It("Relogins if fails refreshing token before attempting to listapps", func() {
				fakeApi.MaxTokenRefresh = 0
				fakeApi.TokenExpiresIn = 12 // Set it to > 10s or Token.Valid() will say false

				previousTokenCounter := fakeApi.TokenCounter

				_, err := http.Get(url.String())
				Expect(err).To(Succeed())
				Expect(fakeApi.TokenCounter).To(Equal(previousTokenCounter + 2))
				Expect(fakeApi.TokenRefreshCounter).To(Equal(1))
			})
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
