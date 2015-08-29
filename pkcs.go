package pbc

import (
	"fmt"
	"os"
	"os/exec"
)

func signPassbook(sign SignConfig, manifestPath, signaturePath string) error {

	args := []string{
		"smime", "-binary", "-sign",
		"-in", manifestPath,
		"-out", signaturePath,
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

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func createCertKey(p12path, cert, key string) error {
	mkCert := exec.Command("openssl", "pkcs12", "-in", p12path, "-clcerts", "-nokeys", "-out", cert)
	if err := mkCert.Run(); err != nil {
		return err
	}

	mkKey := exec.Command("openssl", "pkcs12", "-in", p12path, "-nocerts", "-out", key)

	if err := mkKey.Run(); err != nil {
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
