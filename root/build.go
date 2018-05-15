package root

import (
	"context"
	"io"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type buildOptions struct {
	ImageName   string
	Dockerfile  string
	ContextDir  string
	ProgressBuf io.Writer
	BuildBuf    io.Writer
	BuildArgs   map[string]*string
}

func build(ctx context.Context, cli client.Client, opts *buildOptions) error {
	logrus.Debugf("Running docker build: context: %s, dockerfile: %s", opts.ContextDir, opts.Dockerfile)

	imageBuildOpts := types.ImageBuildOptions{
		Tags:       []string{opts.ImageName},
		Dockerfile: opts.Dockerfile,
		BuildArgs:  opts.BuildArgs,
	}

	buildCtx, buildCtxWriter := io.Pipe()
	go func() {
		err := docker.CreateDockerTarContext(buildCtxWriter, opts.Dockerfile, opts.ContextDir)
		if err != nil {
			buildCtxWriter.CloseWithError(errors.Wrap(err, "creating docker context"))
			return
		}
		buildCtxWriter.Close()
	}()

	progressOutput := streamformatter.NewProgressOutput(opts.ProgressBuf)
	body := progress.NewProgressReader(buildCtx, progressOutput, 0, "", "Sending build context to Docker daemon")

	resp, err := cli.ImageBuild(ctx, body, imageBuildOpts)
	if err != nil {
		return errors.Wrap(err, "docker build")
	}
	defer resp.Body.Close()
	return streamDockerMessages(opts.BuildBuf, resp.Body)
}
