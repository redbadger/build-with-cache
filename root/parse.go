package root

import (
	"os"
	"strconv"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

// parse the Dockerfile into an array of stages
func parse(fileName string) (stages []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	ast, err := parser.Parse(file)
	ss, _, err := instructions.Parse(ast.AST)
	if err != nil {
		return
	}

	stage := 0
	for _, x := range ss {
		name := x.Name
		if name == "" {
			name = strconv.Itoa(stage)
		}
		stages = append(stages, name)
		stage++
	}
	return
}
