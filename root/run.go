package root

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

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
