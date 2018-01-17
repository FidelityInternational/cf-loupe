package applist

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FidelityInternational/cf-loupe/cf"
	gocf "github.com/cloudfoundry-community/go-cfclient"
	version "github.com/hashicorp/go-version"
)

const staleAppMinAge = time.Hour * 24 * 14 // two weeks
const buildpackFreshnessCap = 1            // can be at most 1 version out of date

type AppData struct {
	Apps    []App
	Summary Summary
}

// App contains app information and its buildpack
type App struct {
	Name       string
	UpdatedAt  string
	Buildpack  Buildpack
	IsStale    bool
	Foundation string
	Org        string
	Space      string
	Instances  int
	MemoryMB   int
	State      string
}

// IsHappy returns true if the app is neither stale nor deprecated
func (app App) IsHappy() bool {
	if !app.IsStale && !app.Buildpack.IsDeprecated {
		return true
	}
	return false
}

// Buildpack contains buildpack name and version
type Buildpack struct {
	Name         string
	Version      string
	Freshness    int // corresponds to the number of versions out of date (0 = latest)
	IsDeprecated bool
}

type Summary struct {
	TotalApps      int
	StaleApps      int
	DeprecatedApps int
}

// BuildAppData returns App Data
func BuildAppData(cfClients map[string]cf.IClient, now time.Time) (AppData, error) {
	allApps := []App{}

	foundations, err := getFoundationsAsync(cfClients)
	if err != nil {
		return AppData{}, err
	}

	for foundationName, foundation := range foundations {
		appsForFoundation, err := BuildAppList(foundation, now, foundationName)
		if err != nil {
			return AppData{}, err
		}

		allApps = append(allApps, appsForFoundation...)
	}

	summary := BuildSummary(allApps)

	return AppData{
		Apps:    allApps,
		Summary: summary,
	}, nil
}

// BuildSummary returns a new summary of an app list
func BuildSummary(apps []App) Summary {
	staleApps := 0
	deprecatedApps := 0

	for _, app := range apps {
		if app.IsStale {
			staleApps++
		}
		if app.Buildpack.IsDeprecated {
			deprecatedApps++
		}
	}

	return Summary{
		TotalApps:      len(apps),
		StaleApps:      staleApps,
		DeprecatedApps: deprecatedApps,
	}
}

// SupportStatus is the inverse of deprecation status
func (buildpack Buildpack) SupportStatus() string {
	if buildpack.Version == "" {
		return "no"
	}
	if buildpack.Freshness <= buildpackFreshnessCap {
		return "yes"
	}
	return "no"
}

// BuildAppList takes a list of apps from the go cfclient and a map buildpackGUID to buildpack information, and it returns a list of apps
func BuildAppList(foundation Foundation, now time.Time, foundationName string) ([]App, error) {
	buildpacksMap, err := generateBuildpacks(foundation.GoCFBuildpacks)
	if err != nil {
		return nil, err
	}

	apps := []App{}
	for _, cfClientApp := range foundation.GoCFApps {
		if cfClientApp.UpdatedAt == "" {
			cfClientApp.UpdatedAt = cfClientApp.CreatedAt
		}
		updatedAt, err := time.Parse(time.RFC3339, cfClientApp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		buildpackGUID := cfClientApp.DetectedBuildpackGuid
		var buildpack Buildpack
		if buildpackGUID == "" {
			if cfClientApp.Buildpack == "" {
				buildpack = Buildpack{
					Name:         "Undetected - app unable to start",
					Version:      "Not applicable",
					Freshness:    99,
					IsDeprecated: true,
				}
			} else {
				buildpack = Buildpack{
					Name:         cfClientApp.Buildpack,
					IsDeprecated: true,
				}
			}
		} else {
			var ok bool
			buildpack, ok = buildpacksMap[buildpackGUID]
			if !ok {
				buildpack = Buildpack{
					Name:         "Deleted",
					Freshness:    99,
					IsDeprecated: true,
				}
			}
		}

		spaceGUID := cfClientApp.SpaceGuid
		space, ok := foundation.GoCFSpaces[spaceGUID]
		if !ok {
			continue
		}
		spaceName := space.Name

		orgGUID := space.OrganizationGuid
		org, ok := foundation.GoCFOrgs[orgGUID]
		if !ok {
			continue
		}
		orgName := org.Name

		isStale := now.Sub(updatedAt) >= staleAppMinAge

		apps = append(apps, App{
			Name:       cfClientApp.Name,
			UpdatedAt:  updatedAt.Format("2006-01-02"),
			Buildpack:  buildpack,
			IsStale:    isStale,
			Foundation: foundationName,
			Org:        orgName,
			Space:      spaceName,
			Instances:  cfClientApp.Instances,
			MemoryMB:   cfClientApp.Memory,
			State:      strings.ToLower(cfClientApp.State),
		})
	}
	return apps, nil
}

func generateBuildpacks(gocfbuildpacksMap map[string]gocf.Buildpack) (map[string]Buildpack, error) {
	buildpacksMap := map[string]Buildpack{}

	for guid, gocfBuildpack := range gocfbuildpacksMap {
		name, version, err := parseBuildpackFilename(gocfBuildpack.Filename)
		if err != nil {
			return nil, err
		}
		buildpack := Buildpack{
			Name:      name,
			Version:   version,
			Freshness: 0,
		}
		buildpacksMap[guid] = buildpack
	}

	// generate versions map
	buildpackVersions, err := generateBuildpackVersions(buildpacksMap)
	if err != nil {
		return nil, err
	}

	// update buildpacks map with freshness values
	for guid, buildpack := range buildpacksMap {
		availableVersions := buildpackVersions[buildpack.Name] // 4 => [4.6.1, 4.6.0], 3 => [3.18.0, 3.19.0, 3.20.0]

		// ensure that version is formated according to version library
		parsedVersion, _ := version.NewVersion(buildpack.Version)
		version := parsedVersion.String()

		for _, minorVersions := range availableVersions {
			for index, possibleVersion := range minorVersions {
				if version == possibleVersion {
					freshness := len(minorVersions) - index - 1
					buildpack.Freshness = freshness
					buildpacksMap[guid] = buildpack
					break
				}
			}
		}
	}

	// Update buildpack map with deprecation status
	// Any buildpack two or more versions out of date is considered deprecated
	// Custom buildpacks have an unknown deprecation status, since we can't extract the version
	for guid, buildpack := range buildpacksMap {
		if buildpack.Version != "" && buildpack.Freshness <= buildpackFreshnessCap {
			buildpack.IsDeprecated = false
		} else {
			buildpack.IsDeprecated = true
		}
		buildpacksMap[guid] = buildpack
	}

	return buildpacksMap, nil
}

func generateBuildpackVersions(buildpacksMap map[string]Buildpack) (map[string]map[int][]string, error) {
	// maps buildpack name -> list of versions
	versionsMap := map[string][]string{}
	majorVersionsMap := map[string]map[int][]string{}

	// if no versions array for buildpack, initialize a new versions array
	for _, buildpack := range buildpacksMap {
		_, ok := versionsMap[buildpack.Name]
		if !ok {
			versionsMap[buildpack.Name] = []string{}
		}
	}

	// Add versions to versions array for each buildpack
	for _, buildpack := range buildpacksMap {
		versions := versionsMap[buildpack.Name]

		versions = append(versions, buildpack.Version)
		versionsMap[buildpack.Name] = versions
	}

	// De-duplicate versions array for each buildpack
	for _, buildpack := range buildpacksMap {
		versions := versionsMap[buildpack.Name]
		versionsMap[buildpack.Name] = deduplicateStringArray(versions)
	}

	// Sort versions array according to semantic versioning
	for _, buildpack := range buildpacksMap {
		versions := versionsMap[buildpack.Name]
		majorVersionSplit, err := separateMajorVersions(versions)
		if err != nil {
			return nil, err
		}
		for majorVersion, minorVersions := range majorVersionSplit {
			sortedVersions, err := semverSort(minorVersions)
			if err != nil {
				return nil, err
			}
			majorVersionSplit[majorVersion] = sortedVersions
		}

		majorVersionsMap[buildpack.Name] = majorVersionSplit
	}

	return majorVersionsMap, nil
}

func separateMajorVersions(versions []string) (map[int][]string, error) {
	versionsByMajor := map[int][]string{}

	re, err := regexp.Compile(`^(\d+)?.`)
	if err != nil {
		return nil, err
	}

	for _, version := range versions {
		matches := re.FindStringSubmatch(version)
		if len(matches) > 0 {
			majorVersion, _ := strconv.Atoi(matches[1])
			versionsByMajor[majorVersion] = append(versionsByMajor[majorVersion], version)
		}
	}
	return versionsByMajor, nil
}

// Sorts an array of strings by comparing semantic version
func semverSort(versionsStr []string) ([]string, error) {
	// convert string array into array of version.Version
	versions := make([]*version.Version, len(versionsStr))
	for i, str := range versionsStr {
		v, err := version.NewVersion(str)
		if err != nil {
			return nil, err
		}
		versions[i] = v
	}

	// sort the version.Version array
	sort.Sort(version.Collection(versions))

	// convert version.Version array back into string array
	sortedVersionsStr := make([]string, len(versions))
	for i, v := range versions {
		str := v.String()
		sortedVersionsStr[i] = str
	}

	return sortedVersionsStr, nil
}

func deduplicateStringArray(array []string) []string {
	strMap := map[string]interface{}{}

	for _, elem := range array {
		strMap[elem] = nil
	}

	uniqueElements := []string{}
	for key := range strMap {
		uniqueElements = append(uniqueElements, key)
	}

	return uniqueElements
}

func parseBuildpackFilename(filename string) (string, string, error) {
	var (
		version string
		name    string
	)
	// format:
	// NAME_buildpack-cached-v1.2.3.zip
	// NAME_buildpack-offline-v1.2.3.zip
	// NAME-buildpack-v1.2.3.zip
	// where '-cached' or '-offline' are optional, and 'buildpack' can be preceeded with
	// '_' or '-'
	if filename != "" {
		re, err := regexp.Compile(`^([a-z\-]+)[_-]buildpack(-cached|-offline)?-v([0-9\.]+)\.zip$`)
		if err != nil {
			return "", "", err
		}
		captures := re.FindStringSubmatch(filename)
		if len(captures) < 3 || len(captures) > 4 {
			return parseSpecialJavaBuildpackFilename(filename)
		}

		// name is the second capture since the first is the entire match
		name = captures[1]
		// version will be the last capture
		version = captures[len(captures)-1]
	}

	return name, version, nil
}

func parseSpecialJavaBuildpackFilename(filename string) (string, string, error) {
	re, err := regexp.Compile(`^([a-z]+)[_-]buildpack(-cached|-offline)?-v([0-9_]+)-[a-zA-Z]+-[a-f0-9]{7}\.zip$`)
	if err != nil {
		return "", "", err
	}
	captures := re.FindStringSubmatch(filename)
	if len(captures) < 3 || len(captures) > 4 {
		return "", "", fmt.Errorf("Couldn't parse buildpack filename: \"%s\"", filename)
	}
	name := captures[1]
	version := strings.Replace(captures[len(captures)-1], "_", ".", -1)

	return name, version, nil
}
