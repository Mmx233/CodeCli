package config

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

func TestWindowsReplacementSharingFailurePreservesOriginal(t *testing.T) {
	name := filepath.Join(t.TempDir(), "config")
	original := []byte("default: {username: original}\n")
	if err := os.WriteFile(name, original, 0600); err != nil {
		t.Fatal(err)
	}
	p, err := windows.UTF16PtrFromString(name)
	if err != nil {
		t.Fatal(err)
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer windows.CloseHandle(h)
	v := "new"
	if err = (Store{Path: name}).Set("default.username", &v); err == nil {
		t.Fatal("replacement of a sharing-locked file should fail")
	}
	got, err := os.ReadFile(name)
	if err != nil || string(got) != string(original) {
		t.Fatal("old configuration lost", string(got), err)
	}
}
