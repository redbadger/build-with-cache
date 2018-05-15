package root

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func pull(ctx context.Context, cli *client.Client, ref string) (err error) {
	rc, err := cli.ImagePull(ctx, ref, types.ImagePullOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "pulling image from repository")
	}
	defer rc.Close()
	return streamDockerMessages(os.Stdout, rc)
}
