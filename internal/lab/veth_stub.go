//go:build !linux

package lab

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func VethLink(args []string) error {
	return fmt.Errorf("veth-link runs only inside the lst-lab Docker container")
}

const dockerfileLab = `
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o lst .
FROM alpine:3.21
RUN apk add --no-cache iproute2
COPY --from=build /src/lst /usr/local/bin/lst
ENTRYPOINT ["lst", "lab", "veth-link"]
`

const imageLab = "lst-lab"

func ensureLabImage() error {
	if _, err := exec.Command("docker", "inspect", "--type=image", imageLab).CombinedOutput(); err == nil {
		return nil
	}

	fmt.Fprintln(os.Stderr, "   building lst-lab image...")

	dir, err := os.MkdirTemp("", "lst-lab-build")
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(dir+"/Dockerfile", []byte(dockerfileLab), 0644); err != nil {
		return err
	}

	modRoot := moduleRoot()
	var buf bytes.Buffer
	cmd := exec.Command("docker", "build", "-t", imageLab, "-f", dir+"/Dockerfile", modRoot)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}
	return nil
}

func moduleRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "."
	}
	modPath := strings.TrimSpace(string(out))
	if modPath == "" || modPath == os.DevNull {
		return "."
	}
	dir := modPath[:len(modPath)-len("/go.mod")]
	if dir == "" {
		return "."
	}
	return dir
}

func vethRun(idx int, ctrA, ctrB, ifA, ifB string) error {
	if err := ensureLabImage(); err != nil {
		return fmt.Errorf("lab image: %w", err)
	}

	out, err := exec.Command("docker", "run", "--rm",
		"--privileged", "--pid=host",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		imageLab,
		strconv.Itoa(idx), ctrA, ctrB, ifA, ifB,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}
