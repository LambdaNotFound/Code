package solid_coding

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

/**
 * 8. String to Integer (atoi)
 */
func myAtoi(s string) int {
    s = strings.TrimSpace(s)
    multiplier := 1
    if len(s) == 0 {
        return 0
    } else if s[0] == '-' {
        multiplier = -1
        s = s[1:]
    } else if s[0] == '+' {
        s = s[1:]
    }

    res := 0
    for _, r := range s {
        if !unicode.IsDigit(r) {
            break
        }
        curr, _ := strconv.Atoi(string(r))

        if multiplier == 1 && (res*10 > math.MaxInt32-curr) {
            return math.MaxInt32
        }
        if multiplier == -1 && (-res*10 < math.MinInt32+curr) {
            return math.MinInt32
        }

        res = res*10 + curr
    }
    return multiplier * res
}
