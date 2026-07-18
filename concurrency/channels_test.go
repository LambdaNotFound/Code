package concurrency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnbufferedChannel(t *testing.T) {
	ch := make(chan string) // unbuffered

	go func() {
		ch <- "hello from goroutine" // blocks until main receives
	}()

	msg := <-ch // blocks until the goroutine sends

	assert.Equal(t, "hello from goroutine", msg)
}
