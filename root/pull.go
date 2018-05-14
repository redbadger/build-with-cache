package root

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// pull an image
func pull(ref string) (err error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return
	}
	rc, err := cli.ImagePull(ctx, ref, types.ImagePullOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "pulling image from repository")
	}
	defer rc.Close()
	return streamDockerMessages(os.Stdout, rc)
}
