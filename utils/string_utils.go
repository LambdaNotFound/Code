package utils

import "fmt"

type String string

func (str String) foo() string {
    for i := 0; i < len(str); i++ {
        b := str[i] // byte (uint8) everything’s ASCII, so 1 byte = 1 rune.
        fmt.Printf("s[%d] = %c (byte value: %d)\n", i, b, b)
    }

    for i, r := range str {
        fmt.Printf("index %d: rune %c (Unicode: %U)\n", i, r, r)
    }

	return string(str) + " method on custom type"
}

func substr(str string, start, end int) string {
	runes := []rune(str) // convert to rune slice
	return string(runes[start:end])
}
