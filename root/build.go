package root

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func build(dir, file, tag, flags string, images map[string]string) (out string, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer
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
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	if dir == "-" {
		cmd.Stdin = os.Stdin
	}
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start()
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	fmt.Println()
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	out = string(stdoutBuf.Bytes())
	errStr := string(stderrBuf.Bytes())
	if errStr != "" {
		return out, errors.New(errStr)
	}
	return
}
