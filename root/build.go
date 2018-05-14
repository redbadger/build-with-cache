package root

import (
	"os"
	"os/exec"
)

func build(dir, file, tag string) (out []byte, err error) {
	args := []string{"build", dir, "-f", file}
	if tag != "" {
		args = append(args, "-t", tag)
	}

	build := exec.Command("docker", args...)
	if dir == "-" {
		build.Stdin = os.Stdin
	}
	out, err = build.CombinedOutput()
	return
}
