package applist

import (
	"github.com/FidelityInternational/cf-loupe/cf"
	gocf "github.com/cloudfoundry-community/go-cfclient"
)

type Foundation struct {
	GoCFApps       []gocf.App
	GoCFBuildpacks map[string]gocf.Buildpack
	GoCFOrgs       map[string]gocf.Org
	GoCFSpaces     map[string]gocf.Space
}

type cfClientAppsElement struct {
	cfClientApps []gocf.App
	foundation   string
	err          error
}

type buildpacksMapsElement struct {
	buildpacksMap map[string]gocf.Buildpack
	foundation    string
	err           error
}

type orgMapElement struct {
	orgMap     map[string]gocf.Org
	foundation string
	err        error
}

type spaceMapElement struct {
	spaceMap   map[string]gocf.Space
	foundation string
	err        error
}

func getFoundationsAsync(cfClients map[string]cf.IClient) (map[string]Foundation, error) {
	foundations := map[string]Foundation{}
	for foundationName := range cfClients {
		foundations[foundationName] = Foundation{}
	}

	// channel of app lists for each foundation
	cfClientAppsChannel := make(chan cfClientAppsElement)

	// channel of buildpack maps for each foundation
	buildpacksMapsChannel := make(chan buildpacksMapsElement)

	// channel of org maps for each foundation
	orgMapChannel := make(chan orgMapElement)

	// channel of org maps for each foundation
	spaceMapChannel := make(chan spaceMapElement)

	// Asynchronously fetch the list of apps for each foundation and the map of buildpacks
	for foundation, cfClient := range cfClients {
		err := cfClient.ReAuth()
		if err != nil {
			return foundations, err
		}

		go listAppsAsync(foundation, cfClient, cfClientAppsChannel)
		go getBuildpacksAsync(foundation, cfClient, buildpacksMapsChannel)
		go getOrgsAsync(foundation, cfClient, orgMapChannel)
		go getSpacesAsync(foundation, cfClient, spaceMapChannel)
	}

	// Wait until a list of apps has been fetched from each foundation
	for i := 0; i < len(cfClients); i++ {
		cfClientAppsElem := <-cfClientAppsChannel
		if cfClientAppsElem.err != nil {
			return nil, cfClientAppsElem.err
		}
		foundation := foundations[cfClientAppsElem.foundation]
		foundation.GoCFApps = cfClientAppsElem.cfClientApps
		foundations[cfClientAppsElem.foundation] = foundation
	}
	close(cfClientAppsChannel)

	// Wait until a map of buildpacks has been fetched from each foundation
	for i := 0; i < len(cfClients); i++ {
		buildpacksMapsElem := <-buildpacksMapsChannel
		if buildpacksMapsElem.err != nil {
			return nil, buildpacksMapsElem.err
		}
		foundation := foundations[buildpacksMapsElem.foundation]
		foundation.GoCFBuildpacks = buildpacksMapsElem.buildpacksMap
		foundations[buildpacksMapsElem.foundation] = foundation
	}
	close(buildpacksMapsChannel)

	// Wait until a map of orgs has been fetched from each foundation
	for i := 0; i < len(cfClients); i++ {
		orgMapElem := <-orgMapChannel
		if orgMapElem.err != nil {
			return nil, orgMapElem.err
		}
		foundation := foundations[orgMapElem.foundation]
		foundation.GoCFOrgs = orgMapElem.orgMap
		foundations[orgMapElem.foundation] = foundation
	}
	close(orgMapChannel)

	// Wait until a map of spaces has been fetched from each foundation
	for i := 0; i < len(cfClients); i++ {
		spaceMapElem := <-spaceMapChannel
		if spaceMapElem.err != nil {
			return nil, spaceMapElem.err
		}
		foundation := foundations[spaceMapElem.foundation]
		foundation.GoCFSpaces = spaceMapElem.spaceMap
		foundations[spaceMapElem.foundation] = foundation
	}
	close(spaceMapChannel)

	return foundations, nil
}

func listAppsAsync(foundation string, cfClient cf.IClient, cfClientAppsChannel chan cfClientAppsElement) {
	cfClientApps, err := cfClient.ListApps()
	cfClientAppsChannel <- cfClientAppsElement{
		cfClientApps: cfClientApps,
		foundation:   foundation,
		err:          err,
	}
}

func getBuildpacksAsync(foundation string, cfClient cf.IClient, buildpacksMapsChannel chan buildpacksMapsElement) {
	buildpacksMap, err := cfClient.GetBuildpacks()
	buildpacksMapsChannel <- buildpacksMapsElement{
		buildpacksMap: buildpacksMap,
		foundation:    foundation,
		err:           err,
	}
}

func getOrgsAsync(foundation string, cfClient cf.IClient, orgMapChannel chan orgMapElement) {
	orgMap, err := cfClient.GetOrgs()
	orgMapChannel <- orgMapElement{
		orgMap:     orgMap,
		foundation: foundation,
		err:        err,
	}
}

func getSpacesAsync(foundation string, cfClient cf.IClient, spaceMapChannel chan spaceMapElement) {
	spaceMap, err := cfClient.GetSpaces()
	spaceMapChannel <- spaceMapElement{
		spaceMap:   spaceMap,
		foundation: foundation,
		err:        err,
	}
}
