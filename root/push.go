package root

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func push(ctx context.Context, cli *client.Client, ref string) (err error) {
	auth, err := getEncodedCredentials(ref)
	if err != nil {
		return errors.Wrap(err, "getting credentials")
	}
	rc, err := cli.ImagePush(ctx, ref, types.ImagePushOptions{RegistryAuth: auth})
	if err != nil {
		return errors.Wrap(err, "pushing image to repository")
	}
	defer rc.Close()
	return streamDockerMessages(os.Stdout, rc)
}
