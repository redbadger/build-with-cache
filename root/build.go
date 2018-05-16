package root

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func build(dir, file, tag, flags string, images map[string]string) (out string, err error) {
	var stdoutBuf bytes.Buffer

	args := []string{"build", dir, "-f", file}
	if tag != "" {
		args = append(args, "-t", tag)
	}
	for _, image := range images {
		args = append(args, "--cache-from", image)
	}
	for _, field := range strings.Fields(flags) {
		args = append(args, field)
	}

	fmt.Printf("\nCommand: docker %s\n", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)

	cmd.Stderr = os.Stderr
	if dir == "-" {
		cmd.Stdin = os.Stdin
	}

	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)

	err = cmd.Start()
	if err != nil {
		log.Fatalf("failed to start cmd (%v)", err)
	}
	_, err = io.Copy(stdout, stdoutIn)
	if err != nil {
		log.Fatalf("failed to capture stdout (%v)\n", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Wait() failed with %s\n", err)
	}

	out = string(stdoutBuf.Bytes())

	fmt.Println()
	return
}
