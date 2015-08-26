package pbc

import (
	"os/exec"
)

func signPassbook(sign SignConfig) error {
	cmd := exec.Command("openssl", "smime", "-binary", "-sign",
		"-certfile", sign.Cert,
		"-signer", sign.Signer,
		"-inkey", sign.Key,
		"-in", "manifest.json",
		"-out", "signature",
		"-passin", sign.Pass)
	return cmd.Run()
}

type SignConfig struct {
	Cert   string
	Signer string
	Key    string
	Pass   string
}
