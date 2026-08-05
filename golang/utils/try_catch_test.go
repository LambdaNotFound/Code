package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * golang patterns
 *
 * Go has no try/catch — a panic unwinds the stack until a deferred function
 * calls recover(). recover() only stops the panic when called directly
 * inside a deferred function; it returns nil if there was no panic in
 * progress, otherwise the value passed to panic().
 */

// Try emulates try/catch: it runs fn and converts any panic into an error.
func Try(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("recovered panic: %v", r)
			}
		}
	}()

	fn()
	return nil
}

func Test_recover_catches_panic(t *testing.T) {
	err := Try(func() {
		panic("boom")
	})

	assert.EqualError(t, err, "recovered panic: boom")
}

func Test_recover_preserves_error_panic(t *testing.T) {
	sentinel := errors.New("sentinel failure")

	err := Try(func() {
		panic(sentinel)
	})

	assert.ErrorIs(t, err, sentinel)
}

func Test_recover_no_panic_returns_nil(t *testing.T) {
	err := Try(func() {
		_ = 1 + 1
	})

	assert.NoError(t, err)
}

func Test_recover_must_be_called_directly_in_deferred_func(t *testing.T) {
	// recover() only catches when invoked directly by the deferred function;
	// calling it one level deeper (e.g. via a helper) would not stop the panic.
	result := 0

	func() {
		defer func() {
			recover()
			result = 1
		}()

		panic("boom")
	}()

	assert.Equal(t, 1, result)
}

func Test_recover_then_rethrow(t *testing.T) {
	assert.Panics(t, func() {
		defer func() {
			if r := recover(); r != nil {
				panic(r) // inspect/log, then rethrow
			}
		}()

		panic("original")
	})
}

func Test_finally_style_cleanup_runs_on_panic(t *testing.T) {
	cleanedUp := false

	func() {
		defer func() {
			cleanedUp = true // runs whether or not a panic occurred, like "finally"
			recover()
		}()

		panic("boom")
	}()

	assert.True(t, cleanedUp)
}
