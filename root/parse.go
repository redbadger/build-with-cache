package root

import (
	"os"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
)

// Parse the Dockerfile into an array of stages
func Parse(fileName string) (targets []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	ast, err := parser.Parse(file)
	stages, _, err := instructions.Parse(ast.AST)
	if err != nil {
		return
	}

	for _, x := range stages {
		targets = append(targets, x.Name)
	}
	return
}
