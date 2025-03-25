package utils

type String string

func (str String) foo() string {
	return string(str) + " method on custom type"
}

func substr(str string, start, end int) string {
	runes := []rune(str)
	return string(runes[start:end])
}
