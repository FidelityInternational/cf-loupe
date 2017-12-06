package applist_test

import (
	"time"

	. "github.com/FidelityInternational/cf-loupe/applist"
	gocf "github.com/cloudfoundry-community/go-cfclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build", func() {
	It("returns the correct apps, for all known buildpacks", func() {
		gocfApps := []gocf.App{
			gocf.App{Name: "app-java", DetectedBuildpackGuid: "guid-1", UpdatedAt: "2017-08-01T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-2", UpdatedAt: "2017-08-02T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-3", UpdatedAt: "2017-08-03T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-4", UpdatedAt: "2017-08-04T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-nodejs", DetectedBuildpackGuid: "guid-5", UpdatedAt: "2017-08-05T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-java", DetectedBuildpackGuid: "guid-6", UpdatedAt: "2017-08-06T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-nodejs", DetectedBuildpackGuid: "guid-7", UpdatedAt: "2017-08-07T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-8", UpdatedAt: "2017-08-08T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-9", UpdatedAt: "2017-08-09T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-java", DetectedBuildpackGuid: "guid-10", UpdatedAt: "2017-08-10T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python", DetectedBuildpackGuid: "guid-11", UpdatedAt: "2017-08-11T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-ruby", DetectedBuildpackGuid: "guid-12", UpdatedAt: "2017-08-12T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-dotnet", DetectedBuildpackGuid: "guid-13", UpdatedAt: "2017-08-13T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-ruby", DetectedBuildpackGuid: "guid-14", UpdatedAt: "2017-08-14T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-staticfile_buildpack", DetectedBuildpackGuid: "guid-15", UpdatedAt: "2017-08-15T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-ruby_buildpack", DetectedBuildpackGuid: "guid-16", UpdatedAt: "2017-08-16T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-java_buildpack_offline", DetectedBuildpackGuid: "guid-17", UpdatedAt: "2017-08-17T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-nodejs_buildpack", DetectedBuildpackGuid: "guid-18", UpdatedAt: "2017-08-18T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-go_buildpack", DetectedBuildpackGuid: "guid-19", UpdatedAt: "2017-08-19T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-python_buildpack", DetectedBuildpackGuid: "guid-20", UpdatedAt: "2017-08-20T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-php_buildpack", DetectedBuildpackGuid: "guid-21", UpdatedAt: "2017-08-21T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-binary_buildpack", DetectedBuildpackGuid: "guid-22", UpdatedAt: "2017-08-22T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-dotnet", DetectedBuildpackGuid: "guid-23", UpdatedAt: "2017-08-23T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
			gocf.App{Name: "app-no-buildpack", DetectedBuildpackGuid: "", UpdatedAt: "2017-08-23T12:00:00Z", SpaceGuid: "def456", Instances: 1, Memory: 512, State: "stopped"},
		}

		buildpacksMap := map[string]gocf.Buildpack{
			"guid-1":  gocf.Buildpack{Name: "java-v3_19-company-b7c2d95", Filename: "java-buildpack-v3_19-company-b7c2d95.zip"},
			"guid-2":  gocf.Buildpack{Name: "python-v1_5_23-company-a169424", Filename: "python_buildpack-cached-v1_5_23-company-a169424.zip"},
			"guid-3":  gocf.Buildpack{Name: "python-v1_5_22-company-6d8603d", Filename: "python_buildpack-cached-v1_5_22-company-6d8603d.zip"},
			"guid-4":  gocf.Buildpack{Name: "python-v1_5_21-company-233b817", Filename: "python_buildpack-cached-v1_5_21-company-233b817.zip"},
			"guid-5":  gocf.Buildpack{Name: "nodejs-1_6_3-company-8f66a52", Filename: "nodejs_buildpack-cached-v1_6_3-company-8f66a52.zip"},
			"guid-6":  gocf.Buildpack{Name: "java-v3_18-company-60c71c6", Filename: "java-buildpack-v3_18-company-60c71c6.zip"},
			"guid-7":  gocf.Buildpack{Name: "nodejs-1_6_2-company-0e20d5b", Filename: "nodejs_buildpack-cached-v1_6_2-company-0e20d5b.zip"},
			"guid-8":  gocf.Buildpack{Name: "python-v1_5_20-company-0db0f5e", Filename: "python_buildpack-cached-v1_5_20-company-0db0f5e.zip"},
			"guid-9":  gocf.Buildpack{Name: "python-v1_5_19-company-1588bd4", Filename: "python_buildpack-cached-v1_5_19-company-1588bd4.zip"},
			"guid-10": gocf.Buildpack{Name: "java-v3_17-company-efe5433", Filename: "java-buildpack-v3_17-company-efe5433.zip"},
			"guid-11": gocf.Buildpack{Name: "python-v1_5_18-company-0bbc4c4", Filename: "python_buildpack-cached-v1_5_18-company-0bbc4c4.zip"},
			"guid-12": gocf.Buildpack{Name: "ruby-1_6_35-company-fb501fe", Filename: "ruby_buildpack-cached-v1_6_35-company-fb501fe.zip"},
			"guid-13": gocf.Buildpack{Name: "dotnet-core-company-v1_0_13", Filename: "dotnet-core_buildpack-cached-v1.0.13.zip"},
			"guid-14": gocf.Buildpack{Name: "ruby-1_6_34-company-20586de", Filename: "ruby_buildpack-cached-v1_6_34-company-20586de.zip"},
			"guid-15": gocf.Buildpack{Name: "staticfile_buildpack", Filename: "staticfile_buildpack-cached-v1.4.6.zip"},
			"guid-16": gocf.Buildpack{Name: "ruby_buildpack", Filename: "ruby_buildpack-cached-v1.6.39.zip"},
			"guid-17": gocf.Buildpack{Name: "java_buildpack_offline", Filename: "java-buildpack-offline-v3.16.zip"},
			"guid-18": gocf.Buildpack{Name: "nodejs_buildpack", Filename: "nodejs_buildpack-cached-v1.5.34.zip"},
			"guid-19": gocf.Buildpack{Name: "go_buildpack", Filename: "go_buildpack-cached-v1.8.2.zip"},
			"guid-20": gocf.Buildpack{Name: "python_buildpack", Filename: "python_buildpack-cached-v1.5.18.zip"},
			"guid-21": gocf.Buildpack{Name: "php_buildpack", Filename: "php_buildpack-cached-v4.3.33.zip"},
			"guid-22": gocf.Buildpack{Name: "binary_buildpack", Filename: "binary-buildpack-v1.0.13.zip"},
			"guid-23": gocf.Buildpack{Name: "dotnet-core-buildpack", Filename: "dotnet-core_buildpack-cached-v1.0.18.zip"},
		}

		foundation := Foundation{
			GoCFApps:       gocfApps,
			GoCFBuildpacks: buildpacksMap,
			GoCFOrgs: map[string]gocf.Org{
				"abc123": gocf.Org{
					Name: "APP1234-project-x",
				},
			},
			GoCFSpaces: map[string]gocf.Space{
				"def456": gocf.Space{
					Name:             "DEV",
					OrganizationGuid: "abc123",
				},
			},
		}

		currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
		foundationName := "dev"
		appList, err := BuildAppList(foundation, currentTime, foundationName)
		Expect(err).To(Succeed())

		Expect(appList[0].Buildpack).To(Equal(Buildpack{Name: "java", Version: "3.19", Freshness: 0, IsDeprecated: false}))
		Expect(appList[1].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.23", Freshness: 0, IsDeprecated: false}))
		Expect(appList[2].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.22", Freshness: 1, IsDeprecated: false}))
		Expect(appList[3].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.21", Freshness: 2, IsDeprecated: true}))
		Expect(appList[4].Buildpack).To(Equal(Buildpack{Name: "nodejs", Version: "1.6.3", Freshness: 0, IsDeprecated: false}))
		Expect(appList[5].Buildpack).To(Equal(Buildpack{Name: "java", Version: "3.18", Freshness: 1, IsDeprecated: false}))
		Expect(appList[6].Buildpack).To(Equal(Buildpack{Name: "nodejs", Version: "1.6.2", Freshness: 1, IsDeprecated: false}))
		Expect(appList[7].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.20", Freshness: 3, IsDeprecated: true}))
		Expect(appList[8].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.19", Freshness: 4, IsDeprecated: true}))
		Expect(appList[9].Buildpack).To(Equal(Buildpack{Name: "java", Version: "3.17", Freshness: 2, IsDeprecated: true}))
		Expect(appList[10].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.18", Freshness: 5, IsDeprecated: true}))
		Expect(appList[11].Buildpack).To(Equal(Buildpack{Name: "ruby", Version: "1.6.35", Freshness: 1, IsDeprecated: false}))
		Expect(appList[12].Buildpack).To(Equal(Buildpack{Name: "dotnet-core", Version: "1.0.13", Freshness: 1, IsDeprecated: false}))
		Expect(appList[13].Buildpack).To(Equal(Buildpack{Name: "ruby", Version: "1.6.34", Freshness: 2, IsDeprecated: true}))
		Expect(appList[14].Buildpack).To(Equal(Buildpack{Name: "staticfile", Version: "1.4.6", Freshness: 0, IsDeprecated: false}))
		Expect(appList[15].Buildpack).To(Equal(Buildpack{Name: "ruby", Version: "1.6.39", Freshness: 0, IsDeprecated: false}))
		Expect(appList[16].Buildpack).To(Equal(Buildpack{Name: "java", Version: "3.16", Freshness: 3, IsDeprecated: true}))
		Expect(appList[17].Buildpack).To(Equal(Buildpack{Name: "nodejs", Version: "1.5.34", Freshness: 2, IsDeprecated: true}))
		Expect(appList[18].Buildpack).To(Equal(Buildpack{Name: "go", Version: "1.8.2", Freshness: 0, IsDeprecated: false}))
		Expect(appList[19].Buildpack).To(Equal(Buildpack{Name: "python", Version: "1.5.18", Freshness: 5, IsDeprecated: true}))
		Expect(appList[20].Buildpack).To(Equal(Buildpack{Name: "php", Version: "4.3.33", Freshness: 0, IsDeprecated: false}))
		Expect(appList[21].Buildpack).To(Equal(Buildpack{Name: "binary", Version: "1.0.13", Freshness: 0, IsDeprecated: false}))
		Expect(appList[22].Buildpack).To(Equal(Buildpack{Name: "dotnet-core", Version: "1.0.18", Freshness: 0, IsDeprecated: false}))
		Expect(appList[23].Buildpack).To(Equal(Buildpack{Name: "Undetected - app unable to start", Version: "Not applicable", Freshness: 99, IsDeprecated: true}))
	})

	Context("when there is a custom buildpack", func() {
		It("returns the correct apps", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "",
						Buildpack:             "https://github.com/cloudfoundry/staticfile-buildpack",
						SpaceGuid:             "def456",
						Instances:             2,
						Memory:                512,
						State:                 "STARTED",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{},
				GoCFOrgs: map[string]gocf.Org{
					"abc123": gocf.Org{
						Name: "APP1234-project-x",
					},
				},
				GoCFSpaces: map[string]gocf.Space{
					"def456": gocf.Space{
						Name:             "DEV",
						OrganizationGuid: "abc123",
					},
				},
			}

			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			appList, err := BuildAppList(foundation, currentTime, "dev")
			Expect(err).To(Succeed())
			Expect(appList).To(HaveLen(1))
			Expect(appList[0].Name).To(Equal("app1"))
			Expect(appList[0].Buildpack.Name).To(Equal("https://github.com/cloudfoundry/staticfile-buildpack"))
			Expect(appList[0].Buildpack.Version).To(BeEmpty())
			Expect(appList[0].Buildpack.IsDeprecated).To(BeTrue())
			Expect(appList[0].UpdatedAt).To(Equal("2017-08-12"))
			Expect(appList[0].IsStale).To(Equal(false))
			Expect(appList[0].Foundation).To(Equal("dev"))
			Expect(appList[0].Org).To(Equal("APP1234-project-x"))
			Expect(appList[0].Space).To(Equal("DEV"))
			Expect(appList[0].Instances).To(Equal(2))
			Expect(appList[0].MemoryMB).To(Equal(512))
			Expect(appList[0].State).To(Equal("started"))
		})
	})

	Context("when the Fidelity Java Buildpack is used", func() {
		It("returns the correct buildpack information", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "def456",
						SpaceGuid:             "def456",
						Instances:             23,
						Memory:                64,
						State:                 "STOPPED",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{
					"def456": gocf.Buildpack{
						Name:     "java_buildpack",
						Filename: "java-buildpack-v1_19-fidelity-abc1234.zip",
					},
				},
				GoCFOrgs: map[string]gocf.Org{
					"abc123": gocf.Org{
						Name: "APP1234-project-x",
					},
				},
				GoCFSpaces: map[string]gocf.Space{
					"def456": gocf.Space{
						Name:             "DEV",
						OrganizationGuid: "abc123",
					},
				},
			}

			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			appList, err := BuildAppList(foundation, currentTime, "dev")
			Expect(err).To(Succeed())
			Expect(appList).To(HaveLen(1))
			Expect(appList[0].Name).To(Equal("app1"))
			Expect(appList[0].Buildpack.Name).To(Equal("java"))
			Expect(appList[0].Buildpack.Version).To(Equal("1.19"))
			Expect(appList[0].Buildpack.IsDeprecated).To(BeFalse())
			Expect(appList[0].UpdatedAt).To(Equal("2017-08-12"))
			Expect(appList[0].IsStale).To(Equal(false))
			Expect(appList[0].Foundation).To(Equal("dev"))
			Expect(appList[0].Org).To(Equal("APP1234-project-x"))
			Expect(appList[0].Space).To(Equal("DEV"))
			Expect(appList[0].Instances).To(Equal(23))
			Expect(appList[0].MemoryMB).To(Equal(64))
			Expect(appList[0].State).To(Equal("stopped"))
		})
	})

	Context("when the buildpack can't be found", func() {
		It("returns a meaningful error message", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "def456",
						SpaceGuid:             "def456",
					},
				},
				GoCFOrgs: map[string]gocf.Org{
					"abc123": gocf.Org{
						Name: "APP1234-project-x",
					},
				},
				GoCFSpaces: map[string]gocf.Space{
					"def456": gocf.Space{
						Name:             "DEV",
						OrganizationGuid: "abc123",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{},
			}

			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			appList, err := BuildAppList(foundation, currentTime, "dev")
			Expect(err).To(Succeed())
			Expect(appList).To(HaveLen(1))
			Expect(appList[0].Name).To(Equal("app1"))
			Expect(appList[0].Buildpack.Name).To(Equal("Deleted"))
			Expect(appList[0].Buildpack.Version).To(Equal(""))
			Expect(appList[0].Buildpack.IsDeprecated).To(BeTrue())
			Expect(appList[0].UpdatedAt).To(Equal("2017-08-12"))
			Expect(appList[0].IsStale).To(Equal(false))
			Expect(appList[0].Foundation).To(Equal("dev"))
			Expect(appList[0].Org).To(Equal("APP1234-project-x"))
			Expect(appList[0].Space).To(Equal("DEV"))
			Expect(appList[0].Instances).To(Equal(0))
			Expect(appList[0].MemoryMB).To(Equal(0))
		})
	})

	Context("when the space can't be found", func() {
		It("returns a meaningful error message", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "def456",
						SpaceGuid:             "def456",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{
					"def456": gocf.Buildpack{
						Name:     "java_buildpack",
						Filename: "java-buildpack-v1_19-fidelity-abc1234.zip",
					},
				},
				GoCFOrgs: map[string]gocf.Org{
					"abc123": gocf.Org{
						Name: "APP1234-project-x",
					},
				},
			}

			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			appList, err := BuildAppList(foundation, currentTime, "dev")
			Expect(appList).To(HaveLen(0))
			Expect(err).To(Succeed())
		})
	})

	Context("when the org can't be found", func() {
		It("returns a meaningful error message", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "def456",
						SpaceGuid:             "def456",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{
					"def456": gocf.Buildpack{
						Name:     "java_buildpack",
						Filename: "java-buildpack-v1_19-fidelity-abc1234.zip",
					},
				},
				GoCFSpaces: map[string]gocf.Space{
					"def456": gocf.Space{
						Name:             "DEV",
						OrganizationGuid: "abc123",
					},
				},
			}

			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			appList, err := BuildAppList(foundation, currentTime, "dev")
			Expect(appList).To(HaveLen(0))
			Expect(err).To(Succeed())
		})
	})

	Context("When the buildpack has an unrecognised filename", func() {
		It("returns a meaningful error message", func() {
			foundation := Foundation{
				GoCFApps: []gocf.App{
					gocf.App{
						Name:                  "app1",
						UpdatedAt:             "2017-08-12T16:41:45Z",
						DetectedBuildpackGuid: "def456",
						SpaceGuid:             "def456",
					},
				},
				GoCFBuildpacks: map[string]gocf.Buildpack{
					"def456": gocf.Buildpack{
						Name:     "java_buildpack",
						Filename: "bleh",
					},
				},
				GoCFOrgs: map[string]gocf.Org{
					"abc123": gocf.Org{
						Name: "APP1234-project-x",
					},
				},
				GoCFSpaces: map[string]gocf.Space{
					"def456": gocf.Space{
						Name:             "DEV",
						OrganizationGuid: "abc123",
					},
				},
			}
			currentTime, _ := time.Parse(time.RFC3339, "2017-08-24T12:00:00Z")
			_, err := BuildAppList(foundation, currentTime, "dev")
			Expect(err).To(MatchError("Couldn't parse buildpack filename: \"bleh\""))
		})
	})
})
