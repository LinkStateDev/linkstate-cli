//go:build windows

package lab

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

func newDockerClient() *dockerClient {
	path := dockerPipe()
	return &dockerClient{&http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return winio.DialPipe(path, nil)
			},
		},
		Timeout: 20 * time.Second,
	}}
}

func dockerPipe() string {
	if v := os.Getenv("DOCKER_HOST"); v != "" {
		if strings.HasPrefix(v, "npipe://") {
			return v[8:]
		}
		if strings.HasPrefix(v, "tcp://") {
			return v[6:]
		}
		return v
	}
	return `\\.\pipe\docker_engine`
}
