package pbc

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type CompileOptions struct {
	NoZip       bool
	NoSignature bool
}

// Compiles a passbock pass at the directory specified
func Compile(path string, sign SignConfig, options CompileOptions, out io.Writer) error {
	fmt.Printf("Compiling directory %v\n", path)
	fmt.Println("Packaging files")

	targets, err := gatherTargets(path)
	manifest, err := makeManifest(path, targets)

	if err != nil {
		return err
	}

	err = writeManifest(manifest, path)
	targets = addSet(addSet(targets, "manifest.json"), "signature")

	if err = signPassbook(sign); err != nil {
		return err
	}

	name := "vr.pkpass"

	return packagePassbook(name, targets, out)
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

func packagePassbook(name string, targets []string, out io.Writer) error {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	passbook := zip.NewWriter(file)
	defer passbook.Close()

	for _, target := range targets {
		fout, err := passbook.Create(target)
		if err != nil {
			return err
		}

		if err := writeIn(fout, target); err != nil {
			return err
		}
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
