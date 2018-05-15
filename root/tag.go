package root

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func tag(ctx context.Context, cli *client.Client, source, target string) (err error) {
	msg := fmt.Sprintf("Tagging image %s as %s", source, target)
	err = cli.ImageTag(ctx, source, target)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	return
}
