package pbc

import (
	"fmt"
	"os"
	"os/exec"
)

func signPassbook(sign SignConfig) error {

	args := []string{
		"smime", "-binary", "-sign",
		"-in", "manifest.json",
		"-out", "signature",
		"-outform", "DER",
	}

	args = appendNonEmpty(args, []string{
		"-certfile", sign.Cert,
	})

	args = appendNonEmpty(args, []string{
		"-signer", sign.Signer,
	})

	args = appendNonEmpty(args, []string{
		"-inkey", sign.Key,
	})

	if sign.Pass != "" {
		args = append(args, []string{
			"--passin", fmt.Sprintf("pass:%v", sign.Pass),
		}...)
	}

	cmd := exec.Command("openssl", args...)

	fmt.Printf("command for signing: %v\n", cmd.Args)

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

func appendNonEmpty(base []string, toAdd []string) []string {
	var nonempty = true

	for _, e := range toAdd {
		if e == "" {
			nonempty = false
		}
	}

	if !nonempty {
		return base
	} else {
		return append(base, toAdd...)
	}
}
