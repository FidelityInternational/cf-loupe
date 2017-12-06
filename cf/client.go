package cf

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	gocf "github.com/cloudfoundry-community/go-cfclient"
)

// IClient is an interface forthe Cloud Foundry API
type IClient interface {
	ReAuth() error
	ListApps() ([]gocf.App, error)
	GetBuildpacks() (map[string]gocf.Buildpack, error)
	GetOrgs() (map[string]gocf.Org, error)
	GetSpaces() (map[string]gocf.Space, error)
}

// Client is the concrete implemnetation of Client
type Client struct {
	gocfClient *gocf.Client
}

// BuildClientsFromEnvironment looks at environment variables then instantiates
// a client for each foundation and returns a map, mapping the foundation name to client
func BuildClientsFromEnvironment(env []string) (map[string]IClient, error) {
	foundationConfigs, err := BuildClientConfigFromEnvironment(env)
	if err != nil {
		return nil, err
	}

	cfClients := map[string]IClient{}

	for foundation, config := range foundationConfigs {
		client, err := gocf.NewClient(&config)
		if err != nil {
			return nil, err
		}
		cfClients[foundation] = &Client{gocfClient: client}
	}

	return cfClients, nil
}

var query url.Values = url.Values{
	"results-per-page": []string{
		"100",
	},
}

// BuildClientConfigFromEnvironment looks at environment variables then creates
// a client configuration for each foundation and returns a map, mapping the foundation name to the config
func BuildClientConfigFromEnvironment(env []string) (map[string]gocf.Config, error) {
	foundationConfigs := map[string]gocf.Config{}
	envMap := mapifyEnv(env)

	for i := 1; ; i++ {
		usernameKey := fmt.Sprintf("CF_USERNAME_%d", i)
		passwordKey := fmt.Sprintf("CF_PASSWORD_%d", i)
		apiKey := fmt.Sprintf("CF_API_%d", i)
		foundationKey := fmt.Sprintf("CF_FOUNDATION_%d", i)
		foundation, hasFoundationKey := envMap[foundationKey]
		if !hasFoundationKey {
			break
		}

		username, hasUsernameKey := envMap[usernameKey]
		if !hasUsernameKey {
			return nil, fmt.Errorf("%s env var not found for %s foundation", usernameKey, foundation)
		}
		password, hasPasswordKey := envMap[passwordKey]
		if !hasPasswordKey {
			return nil, fmt.Errorf("%s env var not found for %s foundation", passwordKey, foundation)
		}
		api, hasAPIKey := envMap[apiKey]
		if !hasAPIKey {
			return nil, fmt.Errorf("%s env var not found for %s foundation", apiKey, foundation)
		}

		configFoundation := gocf.Config{
			Username:   username,
			Password:   password,
			ApiAddress: api,
		}
		foundationConfigs[foundation] = configFoundation
	}

	if len(foundationConfigs) == 0 {
		return nil, errors.New("no foundation environment variables found. CF_USERNAME_1, CF_PASSWORD_1, CF_API_1 and CF_FOUNDATION_1 must be set")
	}

	return foundationConfigs, nil
}

func mapifyEnv(envArray []string) map[string]string {
	envMap := map[string]string{}
	for _, envVar := range envArray {
		parts := strings.SplitN(envVar, "=", 2)
		key := parts[0]
		value := parts[1]
		envMap[key] = value
	}
	return envMap
}

// Reauthenticates the client with the api
func (client *Client) ReAuth() error {
	token, err := client.gocfClient.Config.TokenSource.Token()

	if err == nil && token.Valid() {
		return nil // We are authenticated and the token is valid
	}

	// Try to reauthenticate
	cleanConfig := client.gocfClient.Config
	cleanConfig.HttpClient = nil
	cleanConfig.Token = ""
	cleanConfig.ClientID = ""

	newClient, err := gocf.NewClient(&cleanConfig)
	if err != nil {
		return err
	}

	client.gocfClient = newClient
	return nil
}

// ListApps returns the currently deployed apps
func (client *Client) ListApps() ([]gocf.App, error) {
	return client.gocfClient.ListAppsByQuery(query)
}

// GetBuildpacks returns a map of buildpack GUID to buildpack details
func (client *Client) GetBuildpacks() (map[string]gocf.Buildpack, error) {
	buildpacksList, err := client.gocfClient.ListBuildpacks()
	if err != nil {
		return nil, err
	}

	buildpacksMap := map[string]gocf.Buildpack{}
	for _, buildpack := range buildpacksList {
		buildpacksMap[buildpack.Guid] = buildpack
	}

	return buildpacksMap, nil
}

// GetOrgs returns a map of org GUID to org details
func (client *Client) GetOrgs() (map[string]gocf.Org, error) {
	orgList, err := client.gocfClient.ListOrgsByQuery(query)
	if err != nil {
		return nil, err
	}

	orgMap := map[string]gocf.Org{}
	for _, org := range orgList {
		orgMap[org.Guid] = org
	}

	return orgMap, nil
}

// GetSpaces returns a map of space GUID to space details
func (client *Client) GetSpaces() (map[string]gocf.Space, error) {
	spaceList, err := client.gocfClient.ListSpacesByQuery(query)
	if err != nil {
		return nil, err
	}

	spaceMap := map[string]gocf.Space{}
	for _, space := range spaceList {
		spaceMap[space.Guid] = space
	}

	return spaceMap, nil
}
