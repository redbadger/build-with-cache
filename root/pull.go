package root

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// pull an image
func pull(ref string) (err error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return
	}
	_, err = cli.ImagePull(ctx, ref, types.ImagePullOptions{All: true})
	return
}
