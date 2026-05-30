//go:build !windows

package lab

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func newDockerClient() *dockerClient {
	addr, netw := dockerSocket()
	return &dockerClient{&http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, netw, addr)
			},
		},
		Timeout: 20 * time.Second,
	}}
}

func dockerSocket() (addr, network string) {
	if v := os.Getenv("DOCKER_HOST"); v != "" {
		switch {
		case strings.HasPrefix(v, "unix://"):
			return v[7:], "unix"
		case strings.HasPrefix(v, "tcp://"):
			return v[6:], "tcp"
		}
		return v, "tcp"
	}
	return "/var/run/docker.sock", "unix"
}
