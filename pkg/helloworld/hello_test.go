package helloworld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sayHello(t *testing.T) {
	name := "Bob"
	want := "Hello Bob"

	got := sayHello(name)
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
	assert.Equal(t, want, got, "The two strings should be the same.")
}
