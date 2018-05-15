package root

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

func streamDockerMessages(dst io.Writer, src io.Reader) error {
	fd, _ := term.GetFdInfo(dst)
	return jsonmessage.DisplayJSONMessagesStream(src, dst, fd, false, nil)
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
	// proxy := "http://host.docker.internal:3128"
	err = build(ctx, *cli, &buildOptions{
		ImageName:   tag,
		Dockerfile:  file,
		ContextDir:  contextDir,
		ProgressBuf: os.Stdout,
		BuildBuf:    os.Stdout,
		BuildArgs:   map[string]*string{
			// "http_proxy":  &proxy,
			// "https_proxy": &proxy,
		},
	})
	if err != nil {
		err = fmt.Errorf("Error running docker build: %s", err)
		return
	}
	return
}
