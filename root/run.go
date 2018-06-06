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
func Run(contextDir, file, imgName, cache, flags string) (err error) {
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
	if imgName != "" {
		reader, err = os.Open(file)
		if err != nil {
			return
		}
		stages, images, err = parseDockerfile(reader, imgName, cache)
		if err != nil {
			return
		}
		for _, stage := range stages {
			img := images[stage]
			fmt.Printf("Pulling: %s\n", img)
			err = pull(ctx, cli, img)
			if err != nil {
				fmt.Printf("pulling %s: %s\n", img, err)
			}
		}
	}

	out, err := build(contextDir, file, imgName, flags, images)
	if err != nil {
		err = fmt.Errorf("Error running docker build: %s", err)
		return
	}

	var names map[string]string
	if imgName != "" {
		names, err = parseBuildOutput(out, imgName, cache, stages)
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
