package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func IsValidLuhn(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	sum := 0
	parity := len(number) % 2

	for i, r := range number {
		digit, err := strconv.Atoi(string(r))
		if err != nil {
			return false
		}

		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

func FindModuleRoot(dir string) (string, error) {
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
