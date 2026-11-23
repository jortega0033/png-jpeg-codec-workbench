package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpAndVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"--help"}, &out, &errOut); code != 0 || !strings.Contains(out.String(), "usage: codeclab") {
		t.Fatalf("help: code %d, out %q", code, out.String())
	}
	out.Reset()
	if code := run([]string{"--version"}, &out, &errOut); code != 0 || !strings.HasPrefix(out.String(), "codeclab ") {
		t.Fatalf("version: code %d, out %q", code, out.String())
	}
}

func TestStubExitsNonZero(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"serve"}, &out, &errOut); code == 0 {
		t.Fatal("serve should exit non-zero")
	}
	if got := errOut.String(); got != "not implemented yet: see #2\n" {
		t.Fatalf("stderr %q", got)
	}
}
