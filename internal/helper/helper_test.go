package helper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{"Correct number", "48212146351759", true},
		{"Incorrect number", "4821214635175", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidLuhn(tt.number); got != tt.want {
				t.Errorf("IsValidLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindModuleRoot(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	root, err := FindModuleRoot(wd)
	if err != nil {
		t.Fatalf("FindModuleRoot() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Errorf("go.mod not found in root = %v: %v", root, err)
	}
	if !strings.HasPrefix(wd, root) {
		t.Errorf("root (%v) is not parent of wd (%v)", root, wd)
	}
}
