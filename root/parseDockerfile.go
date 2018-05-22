package root

import (
	"fmt"
	"io"
	"strconv"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

// parseDockerfile the Dockerfile into an array of stages
func parseDockerfile(reader io.Reader, imgName string) (stages []string, images map[string]string, err error) {
	ast, err := parser.Parse(reader)
	ss, _, err := instructions.Parse(ast.AST)
	if err != nil {
		return
	}

	imgNamed, err := reference.ParseNamed(imgName)
	if err != nil {
		return
	}

	stage := 0
	stages = make([]string, len(ss))
	images = make(map[string]string)
	for _, x := range ss {
		name := x.Name
		if name == "" {
			name = strconv.Itoa(stage)
		}
		stages[stage] = name
		images[name] = fmt.Sprintf("%s-%s", reference.TrimNamed(imgNamed), name)
		stage++
	}
	return
}
