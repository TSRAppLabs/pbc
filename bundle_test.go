package pbc

import (
	"testing"
)

func TestClean(t *testing.T) {
	subject := clean("/var/tmp/hello", "/var/tmp")
	if subject != "hello" {
		t.Errorf("expected %v, but got %v", "hello", subject)
	}

}
