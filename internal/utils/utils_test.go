package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetProjectRoot(t *testing.T) {
	root := getProjectRoot()
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Errorf("getProjectRoot() returned %s, but go.mod not found there: %v", root, err)
	}
}

func TestLoadDotEnvFromSubdir(t *testing.T) {
	// Change working directory to a subdirectory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// internal/utils is where this test is
	err := os.Chdir(filepath.Join(getProjectRoot(), "internal", "utils"))
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// This should not panic if it finds the .env
	entries := loadDotEnv()
	if len(entries) == 0 {
		t.Error("loadDotEnv() returned empty entries")
	}
}
