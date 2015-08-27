package pbc

import (
	"fmt"
	"os"
	"os/exec"
)

func signPassbook(sign SignConfig) error {

	baseArgs := []string{
		"smime", "-binary", "-sign",
		"-in", "manifest.json",
		"-outform", "DER",
	}

	args := append(baseArgs, []string{
		"-certfile", sign.Cert,
		"-signer", sign.Signer,
		"-inkey", sign.Key,
	}...)

	if sign.Pass != "" {
		args = append(args, []string{
			"--passin", fmt.Sprintf("pass:%v", sign.Pass),
		}...)
	}

	cmd := exec.Command("openssl", args...)

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
