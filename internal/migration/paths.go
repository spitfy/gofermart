package migration

import (
	"fmt"
	"os"
	"path/filepath"
)

func findModuleRoot(dir string) (string, error) {
	for {
		gomod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(gomod); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found in any parent directory")
}

func getMigrationsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	moduleRoot, err := findModuleRoot(wd)
	if err != nil {
		return "", err
	}
	migrationsPath := filepath.Join(moduleRoot, "migrations")
	return migrationsPath, nil
}
