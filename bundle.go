package pbc

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Compiles a passbock pass at the directory specified
func Compile(path string, profile Profile, out io.Writer) error {
	fmt.Printf("Compiling directory %v\n", path)
	fmt.Println("Packaging files")

	targets, err := gatherTargets(path)
	manifest, err := makeManifest(path, targets)

	if err != nil {
		return err
	}

	manifestPath, err := writeManifest(manifest, path)

	sign := SignConfig{
		Cert:   filepath.Join(filepath.Join(getDataDir(), "wwdr.pem")),
		Key:    profile.getKeyPath(),
		Signer: profile.getCertPath(),
		Pass:   "",
	}

	signFile, err := ioutil.TempFile("", "signature")

	if err != nil {
		return err
	}

	signFile.Close()

	if err = signPassbook(sign, manifestPath, signFile.Name()); err != nil {
		return err
	}

	return packagePassbook(path, targets, manifestPath, signFile.Name(), out)
}

func gatherTargets(path string) ([]string, error) {
	targets, err := findTargets(path)

	if err != nil {
		return nil, err
	}

	for _, target := range targets {
		fmt.Printf("\t%v\n", target)
	}

	return targets, nil
}

func packagePassbook(root string, targets []string, manifestPath, sigPath string, out io.Writer) error {
	passbook := zip.NewWriter(out)
	defer passbook.Close()

	for _, target := range targets {
		fout, err := passbook.Create(clean(target, root))
		if err != nil {
			return err
		}

		if err := writeIn(fout, target); err != nil {
			return err
		}
	}

	manOut, err := passbook.Create("manifest.json")
	if err != nil {
		writeIn(manOut, manifestPath)
	}

	sigOut, err := passbook.Create("signature")
	if err != nil {
		writeIn(sigOut, sigPath)
	}

	return nil
}

func writeIn(fout io.Writer, target string) error {
	fin, err := os.Open(target)

	if err != nil {
		return err
	}

	defer fin.Close()

	io.Copy(fout, fin)
	return nil
}

func clean(target, root string) string {
	targetParts := strings.Split(target, string(filepath.Separator))
	rootParts := strings.Split(root, string(filepath.Separator))

	return filepath.Join(stripRoot(targetParts, rootParts)...)
}

func stripRoot(targetParts, rootParts []string) []string {
	i := 0

	for i < len(targetParts) && i < len(rootParts) && targetParts[i] == rootParts[i] {
		i++
	}

	return targetParts[i:]
}
