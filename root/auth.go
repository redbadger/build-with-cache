package root

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
)

func makeAuthStr() (authStr string, err error) {
	authConfig := types.AuthConfig{
		Username: "admin",
		Password: "admin123",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
