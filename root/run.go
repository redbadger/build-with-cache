package root

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/docker/distribution/reference"
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
func Run(contextDir, file, tag string) (err error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return
	}
	var stages []string
	if tag != "" {
		stages, err = parse(file)
		if err != nil {
			return
		}
		var ref reference.Named
		for _, stage := range stages {
			ref, err = reference.ParseNamed(tag)
			if err != nil {
				return
			}
			img := fmt.Sprintf("%s-%s", reference.TrimNamed(ref), stage)
			fmt.Printf("Pulling: %s\n", img)
			err = pull(ctx, cli, img)
			if err != nil {
				fmt.Printf("pulling %s: %s\n", img, err)
			}
		}
	}
	out, err := build(contextDir, file, tag)
	fmt.Printf("%s\n", out)
	if err != nil {
		err = fmt.Errorf("Error running docker build: %s", err)
		return
	}

	return
}
