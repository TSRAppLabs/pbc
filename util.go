package pbc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func addSet(targets []string, str string) []string {
	found := false
	for _, target := range targets {
		if target == str {
			found = true
		}
	}

	if !found {
		return append(targets, str)
	}
	return targets
}

func replaceTilde(path string) string {
	parts := getParts(path)

	for i, part := range parts {
		if part == "~" {
			parts[i] = "${HOME}"
		}
	}

	return filepath.Join(parts...)
}

func ExpandPath(path string) string {
	replace := []string{}
	for _, v := range os.Environ() {
		key, val := getKeyValue(v)
		replace = append(replace, []string{key, val}...)
	}

	return strings.NewReplacer(replace...).Replace(path)
}

func getKeyValue(val string) (string, string) {
	parts := strings.SplitN(val, "=", 2)

	if len(parts) > 1 {
		return fmt.Sprintf("${%v}", parts[0]), parts[1]
	} else if len(parts) > 0 {
		return fmt.Sprintf("${%v}", parts[0]), ""
	}

	return "", ""
}

func getParts(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}
