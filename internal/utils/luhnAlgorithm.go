package utils

import (
	"strconv"
	"unicode"
)

func LuhnCheck(number string) bool {
	var sum int
	alt := false

	// идём справа налево
	for i := len(number) - 1; i >= 0; i-- {
		r := rune(number[i])
		if !unicode.IsDigit(r) {
			return false // если встретился нецифровой символ
		}
		n, _ := strconv.Atoi(string(r))

		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
