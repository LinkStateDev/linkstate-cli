package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
)

func cdInto(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if err := os.Chdir(abs); err != nil {
		return err
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = abs

	_ = cmd.Run()
	return nil
}
