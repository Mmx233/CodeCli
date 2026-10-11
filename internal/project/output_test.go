package project

import (
	"bytes"
	"strings"
	"testing"
)

func TestCredentialRedactionAcrossWritesAndPromptDelivery(t *testing.T) {
	var out bytes.Buffer
	w := redactedWriter{dst: &out}
	for _, chunk := range []string{"Cloning https://user:pa", "ssword@host/repo?tok", "en=secret\r", "Password for '", "https://user@ho", "st/repo': "} {
		if _, err := w.Write([]byte(chunk)); err != nil {
			t.Fatal(err)
		}
	}
	got := out.String()
	if strings.Contains(got, "password") || strings.Contains(got, "secret") || strings.Contains(got, "token=") || strings.Contains(got, "user@") {
		t.Fatal(got)
	}
	if !strings.HasSuffix(got, "https://host/repo': ") {
		t.Fatal("prompt not delivered before process completion", got)
	}
	if err := w.flush(); err != nil {
		t.Fatal(err)
	}
}
