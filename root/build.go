package root

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func createCommand(args []string) []string {
	parts := []string{"docker", "build"}
	skip := false
	for _, x := range args {
		if skip {
			skip = false
			continue
		}
		if x == "--cache" {
			skip = true
			continue
		}
		if strings.HasPrefix(x, "--cache=") {
			continue
		}
		parts = append(parts, x)
	}
	return parts
}

func build(dir string, images map[string]string) (out string, err error) {
	var stdoutBuf bytes.Buffer

	args := createCommand(os.Args[1:])
	for _, image := range images {
		args = append(args, "--cache-from", image)
	}

	log.WithField("command", strings.Join(args, " ")).Info("Executing")
	cmd := exec.Command(args[0], args[1:]...)

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
