package root

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func push(ctx context.Context, cli *client.Client, ref string) (err error) {
	authConfig := types.AuthConfig{
		Username: "admin",
		Password: "admin123",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	rc, err := cli.ImagePush(ctx, ref, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return errors.Wrap(err, "pushing image to repository")
	}
	defer rc.Close()
	return streamDockerMessages(os.Stdout, rc)
}
