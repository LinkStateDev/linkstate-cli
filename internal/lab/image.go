package lab

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const dockerfileFRR = `
FROM quay.io/frrouting/frr:10.2.6
RUN apk add --no-cache dropbear && \
    sed -i 's/^bgpd=no/bgpd=yes/' /etc/frr/daemons && \
    touch /etc/frr/bgpd.conf && chown frr:frr /etc/frr/bgpd.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
`

const entrypointFRR = `#!/bin/sh
mkdir -p /etc/dropbear /root/.ssh
[ -f /etc/dropbear/dropbear_rsa_host_key ] || \
    dropbearkey -t rsa -f /etc/dropbear/dropbear_rsa_host_key 2>/dev/null
if [ -n "$SSH_PUBKEY" ]; then
    echo "$SSH_PUBKEY" > /root/.ssh/authorized_keys
fi
dropbear -E -p 22 2>/dev/null
exec /usr/lib/frr/docker-start
`

func ensureFRRImage() error {
	if _, err := exec.Command("docker", "inspect", "--type=image", imageFRR).CombinedOutput(); err == nil {
		return nil
	}
	if _, err := exec.Command("docker", "pull", imageFRR).CombinedOutput(); err == nil {
		return nil
	}

	fmt.Fprintln(os.Stderr, "   building lab-frr image...")

	dir, err := os.MkdirTemp("", "lab-frr-build")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(dockerfileFRR), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "entrypoint.sh"), []byte(entrypointFRR), 0755); err != nil {
		return err
	}

	var buf bytes.Buffer
	cmd := exec.Command("docker", "build", "-t", imageFRR, dir)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}
	return nil
}
