package root

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

func getRegistryHostnameFromName(name string) (hostname string, err error) {
	named, err := reference.ParseNamed(name)
	if err != nil {
		return
	}
	hostname, _ = reference.SplitHostname(named)
	return
}

func getEncodedCredentials(ref string) (authStr string, err error) {
	registryHostname, err := getRegistryHostnameFromName(ref)
	configFile := config.LoadDefaultConfigFile(os.Stderr)

	authConfig, err := configFile.GetAuthConfig(registryHostname)
	if err != nil {
		return "", errors.Wrap(err, "reading credentials from store")
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", errors.Wrap(err, "encoding json credentials")
	}

	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
