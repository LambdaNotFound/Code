package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isMatch(t *testing.T) {
	assert.False(t, isMatch("aa", "a"))
	assert.True(t, isMatch("aa", "*"))
	assert.False(t, isMatch("cb", "?a"))
	assert.True(t, isMatch("adceb", "*a*b"))
	assert.False(t, isMatch("acdcb", "a*c?b"))
	assert.True(t, isMatch("", "*"))
	assert.True(t, isMatch("", ""))
	assert.False(t, isMatch("a", ""))
	assert.True(t, isMatch("abc", "a?c"))
	assert.True(t, isMatch("abc", "***"))
}
