package lab

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

var labDir string

func SetLabDir(dir string) {
	labDir = dir
}

func labConfigDir() string {
	if labDir != "" {
		return filepath.Join(labDir, ".ssh")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".linkstate", ".ssh")
}

func sshKeyPath() string  { return filepath.Join(labConfigDir(), "id_rsa") }
func sshPubPath() string  { return filepath.Join(labConfigDir(), "id_rsa.pub") }
func SSHKeyPath() string  { return sshKeyPath() }

func ensureSSHKey() error {
	dir := labConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	if _, err := os.Stat(sshKeyPath()); err == nil {
		return nil
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	block, err := ssh.MarshalPrivateKey(key, "")
	if err != nil {
		return fmt.Errorf("marshal key: %w", err)
	}
	if err := os.WriteFile(sshKeyPath(), pem.EncodeToMemory(block), 0600); err != nil {
		return fmt.Errorf("write key: %w", err)
	}

	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return fmt.Errorf("pubkey: %w", err)
	}
	return os.WriteFile(sshPubPath(), ssh.MarshalAuthorizedKey(pub), 0644)
}

func readPubKey() (string, error) {
	data, err := os.ReadFile(sshPubPath())
	return string(data), err
}
