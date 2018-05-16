package root

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

func streamDockerMessages(dst io.Writer, src io.Reader) error {
	fd, _ := term.GetFdInfo(dst)
	isTerminal := terminal.IsTerminal(int(fd))
	return jsonmessage.DisplayJSONMessagesStream(src, dst, fd, isTerminal, nil)
}

// Run the root command
func Run(contextDir, file, imgTag, flags string) (err error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return
	}
	var (
		stages []string
		images map[string]string
		reader io.Reader
	)
	if imgTag != "" {
		reader, err = os.Open(file)
		if err != nil {
			return
		}
		stages, images, err = parseDockerfile(reader, imgTag)
		for _, stage := range stages {
			img := images[stage]
			fmt.Printf("Pulling: %s\n", img)
			err = pull(ctx, cli, img)
			if err != nil {
				fmt.Printf("pulling %s: %s\n", img, err)
			}
		}
	}

	out, err := build(contextDir, file, imgTag, flags, images)
	if err != nil {
		err = fmt.Errorf("Error running docker build: %s", err)
		return
	}

	var names map[string]string
	if imgTag != "" {
		names, err = parseBuildOutput(out, imgTag, stages)
		if err != nil {
			return
		}
		for sha, name := range names {
			err = tag(ctx, cli, sha, name)
			if err != nil {
				return
			}
			err = push(ctx, cli, name)
			if err != nil {
				return
			}
		}
	}

	return
}
