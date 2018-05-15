package root

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

func build(dir, file, tag string) (out string, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	args := []string{"build", dir, "-f", file}
	if tag != "" {
		args = append(args, "-t", tag)
	}

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
