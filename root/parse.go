package root

import (
	"fmt"
	"os"
	"strconv"

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

	stage := 0
	for _, x := range stages {
		fmt.Printf("%+v", x)
		name := x.Name
		if name == "" {
			name = strconv.Itoa(stage)
		}
		targets = append(targets, name)
		stage++
	}
	return
}
