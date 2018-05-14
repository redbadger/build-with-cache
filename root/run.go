package root

import (
	"fmt"
)

// Run the root command
func Run(context, file, tag string) (err error) {
	var stages []string
	if tag != "" {
		stages, err = parse(file)
		if err != nil {
			return
		}
		for _, stage := range stages {
			img := fmt.Sprintf("%s-%s", tag, stage)
			fmt.Println(img)
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
