package helper

import (
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
