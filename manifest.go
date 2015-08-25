package pbc

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Compiles a passbock pass at the directory specified
func Compile(path string) error {
	fmt.Printf("Compiling directory %v\n", path)
	return createManifest(path)
}

//TODO: make it work with nest filesystem structure
func createManifest(path string) error {
	content := make(map[string]string)

	infos, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	for _, info := range infos {
		name := info.Name()
		if !ignoreInManifest(name) {
			content[name], err = getHashForFile(filepath.Join(path, name))
			if err != nil {
				return err
			}
		}
	}

	writeManifest(content, path)

	return nil
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

/*
  Currently files that should be ignored in the manifest of a passbook.
*/
func ignoreInManifest(name string) bool {
	if name == "manifest.json" {
		return true
	}

	config := GetConfig()

	for _, pattern := range config.IgnorePatterns {
		if match, err := filepath.Match(pattern, name); err == nil && match {
			return true
		}
	}

	return false
}

/*

*/
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
