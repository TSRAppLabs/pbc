package pbc

import (
	"os"
	"os/exec"
	"strings"
)

func signPassbook(sign SignConfig) error {
	cmd := exec.Command("openssl", "smime", "-binary", "-sign",
		"-certfile", sign.Cert,
		"-signer", sign.Signer,
		"-inkey", sign.Key,
		"-in", "manifest.json",
		"-out", "signature",
		"-passin", strings.Join([]string{"pass", sign.Pass}, ":"))

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

type SignConfig struct {
	Cert   string
	Signer string
	Key    string
	Pass   string
}
