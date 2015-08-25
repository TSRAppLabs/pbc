package pbc

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Compiles a passbock pass at the directory specified
func Compile(path string) error {
	fmt.Printf("Compiling directory %v\n", path)
	fmt.Println("Packaging files")

	targets, err := findTargets(path)

	if err != nil {
		return err
	}

	for _, target := range targets {
		fmt.Printf("\t%v\n", target)
	}

	manifest, err := makeManifest(path, targets)

	if err != nil {
		return err
	}

	err = writeManifest(manifest, path)

	return nil
}

func makeManifest(root string, targets []string) (map[string]string, error) {
	content := make(map[string]string)

	for _, target := range targets {
		hash, err := getHashForFile(filepath.Join(root, target))
		if err != nil {
			return nil, err
		}
		content[target] = hash
	}

	delete(content, "manifest.json")

	return content, nil
}

func writeManifest(content map[string]string, root string) error {
	file, err := os.Create(filepath.Join(root, "manifest.json"))
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(content)
	if err != nil {
		return err
	}

	return nil
}

func getHashForFile(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
