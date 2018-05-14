package root

import (
	"fmt"
	"os"
	"os/exec"
)

// Root is the execution of the root command
func Root(dir string) error {
	build := exec.Command("docker", "build", dir)
	if dir == "-" {
		build.Stdin = os.Stdin
	}
	out, err := build.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}
