package pbc

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func makeManifest(root string, targets []string) (map[string]string, error) {
	content := make(map[string]string)

	for _, target := range targets {
		hash, err := getHashForFile(target)
		if err != nil {
			return nil, err
		}
		content[target] = hash
	}

	delete(content, "manifest.json")

	return content, nil
}

func writeManifest(content map[string]string, root string) (string, error) {
	file, err := ioutil.TempFile("", "manifest")
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(content)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
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
