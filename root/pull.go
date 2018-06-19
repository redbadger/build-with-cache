package root

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func pull(ctx context.Context, cli *client.Client, ref string) (err error) {
	auth, err := getEncodedCredentials(ref)
	if err != nil {
		return errors.Wrap(err, "getting credentials")
	}
	rc, err := cli.ImagePull(ctx, ref, types.ImagePullOptions{All: true, RegistryAuth: auth})
	if err != nil {
		return errors.Wrap(err, "pulling image from repository")
	}
	defer rc.Close()
	return streamDockerMessages(os.Stdout, rc)
}
