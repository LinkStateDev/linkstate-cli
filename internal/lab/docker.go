package lab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type dockerClient struct{ *http.Client }

func (d *dockerClient) req(method, path string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, "http://localhost"+path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	return d.Do(r)
}

func (d *dockerClient) json(method, path string, body, dest any) error {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	resp, err := d.req(method, path, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return apiError{path, resp.StatusCode, strings.TrimSpace(string(b))}
	}
	if dest != nil {
		return json.NewDecoder(resp.Body).Decode(dest)
	}
	return nil
}

type apiError struct {
	path string
	code int
	body string
}

func (e apiError) Error() string {
	return fmt.Sprintf("docker %s (%d): %s", e.path, e.code, e.body)
}

func (d *dockerClient) containerExists(name string) bool {
	var v any
	return d.json("GET", "/containers/"+name+"/json", nil, &v) == nil
}

func (d *dockerClient) containerRunning(name string) bool {
	var v struct {
		State struct {
			Running bool `json:"Running"`
		} `json:"State"`
	}
	if err := d.json("GET", "/containers/"+name+"/json", nil, &v); err != nil {
		return false
	}
	return v.State.Running
}

func (d *dockerClient) containerPID(name string) (int, error) {
	var v struct {
		State struct{ Pid int `json:"Pid"` } `json:"State"`
	}
	if err := d.json("GET", "/containers/"+name+"/json", nil, &v); err != nil {
		return 0, err
	}
	return v.State.Pid, nil
}

func (d *dockerClient) createContainer(name, hostname, sshPort, sshPub string) error {
	if d.containerExists(name) {
		d.rmContainer(name)
	}
	hc := map[string]any{"Privileged": true}
	if sshPort != "" {
		hc["PortBindings"] = map[string]any{
			"22/tcp": []map[string]string{{"HostPort": sshPort}},
		}
	}
	body := map[string]any{
		"Image":        imageFRR,
		"Hostname":     hostname,
		"HostConfig":   hc,
		"ExposedPorts": map[string]struct{}{"22/tcp": {}},
		"Labels":       map[string]string{"lst-lab": "true"},
	}
	if sshPub != "" {
		body["Env"] = []string{"SSH_PUBKEY=" + sshPub}
	}
	var v struct{ Id string }
	if err := d.json("POST", "/containers/create?name="+name, body, &v); err != nil {
		return fmt.Errorf("create %s: %w", name, err)
	}
	return d.json("POST", "/containers/"+v.Id+"/start", nil, nil)
}

func (d *dockerClient) rmContainer(name string) {
	d.json("POST", "/containers/"+name+"/stop?t=0", nil, nil)
	d.json("DELETE", "/containers/"+name+"?force=1", nil, nil)
}

func (d *dockerClient) waitContainer(name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if d.containerRunning(name) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %q", name)
}

func (d *dockerClient) exec(name string, cmd ...string) (string, error) {
	var v struct{ Id string }
	if err := d.json("POST", "/containers/"+name+"/exec",
		map[string]any{"Cmd": cmd, "AttachStdout": true, "AttachStderr": true}, &v); err != nil {
		return "", err
	}
	resp, err := d.req("POST", "/exec/"+v.Id+"/start", body(map[string]any{"Detach": false, "Tty": false}))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	drainDemux(resp.Body, &buf)

	var info struct{ ExitCode int `json:"ExitCode"` }
	d.json("GET", "/exec/"+v.Id+"/json", nil, &info)

	out := strings.TrimSpace(buf.String())
	if info.ExitCode != 0 {
		return out, fmt.Errorf("exit %d: %s", info.ExitCode, out)
	}
	return out, nil
}

func (d *dockerClient) execOk(name string, cmd ...string) {
	out, err := d.exec(name, cmd...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "exec %s %v: %v\n%s\n", name, cmd, err, out)
	}
}

func (d *dockerClient) listLabContainers() ([]containerRow, error) {
	var result []containerRow
	err := d.json("GET",
		"/containers/json?all=1&filters="+url.QueryEscape(`{"label":["lst-lab=true"]}`),
		nil, &result)
	return result, err
}

type containerRow struct {
	Names []string `json:"Names"`
	State string   `json:"State"`
}

func body(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func drainDemux(r io.Reader, w io.Writer) {
	hdr := make([]byte, 8)
	for {
		if _, err := io.ReadFull(r, hdr); err != nil {
			break
		}
		n := int(hdr[4])<<24 | int(hdr[5])<<16 | int(hdr[6])<<8 | int(hdr[7])
		if n == 0 {
			continue
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil {
			break
		}
		w.Write(buf)
	}
}
