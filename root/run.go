package root

import (
	"fmt"
	"io"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

func streamDockerMessages(dst io.Writer, src io.Reader) error {
	fd, _ := term.GetFdInfo(dst)
	return jsonmessage.DisplayJSONMessagesStream(src, dst, fd, false, nil)
}

// Run the root command
func Run(context, file, tag string) (err error) {
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
			pull(img)
		}
	}
	out, err := build(context, file, tag)
	fmt.Printf("%s\n", out)
	if err != nil {
		err = fmt.Errorf("Error running docker build: %s", err)
		return
	}
	return
}
