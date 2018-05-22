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
func parseDockerfile(reader io.Reader, imgName, cache string) (stages []string, imgNames map[string]string, err error) {
	ast, err := parser.Parse(reader)
	ss, _, err := instructions.Parse(ast.AST)
	if err != nil {
		return
	}

	imgNamed, err := getNamed(imgName, cache)
	if err != nil {
		return
	}
	stage := 0
	stages = make([]string, len(ss))
	imgNames = make(map[string]string)
	for _, s := range ss {
		stageName := s.Name
		if stageName == "" {
			stageName = strconv.Itoa(stage)
		}
		stages[stage] = stageName
		named := reference.TrimNamed(imgNamed)
		imgNames[stageName] = fmt.Sprintf("%s-%s", named, stageName)
		stage++
	}
	return
}

func getNamed(imgName, cache string) (imgNamed reference.Named, err error) {
	imgNamed, err = reference.ParseNamed(imgName)
	if err != nil {
		return
	}
	var (
		domain  string
		nameStr string
	)
	domain, nameStr = reference.SplitHostname(imgNamed)
	if cache != "" {
		domain = cache
	}
	imgNamed, err = reference.WithName(domain + "/" + nameStr)
	if err != nil {
		return
	}
	return
}
