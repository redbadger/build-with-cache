package root

import (
	"os"
	"os/exec"
)

// Root is the execution of the root command
func Root(dir, file string) (out []byte, err error) {
	build := exec.Command("docker", "build", dir, "-f", file)
	if dir == "-" {
		build.Stdin = os.Stdin
	}
	out, err = build.CombinedOutput()
	return
}
