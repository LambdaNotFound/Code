package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l := NewLogger()
	assert.True(t, l.ShouldPrintMessage(1, "foo"))
	assert.True(t, l.ShouldPrintMessage(2, "bar"))
	assert.False(t, l.ShouldPrintMessage(3, "foo"))
	assert.False(t, l.ShouldPrintMessage(8, "bar"))
	assert.False(t, l.ShouldPrintMessage(10, "foo"))
	assert.True(t, l.ShouldPrintMessage(11, "foo"))
}

func TestLoggerV2(t *testing.T) {
	l := NewLoggerV2()
	assert.True(t, l.ShouldPrintMessage(1, "foo"))
	assert.True(t, l.ShouldPrintMessage(2, "bar"))
	assert.False(t, l.ShouldPrintMessage(3, "foo"))
	assert.False(t, l.ShouldPrintMessage(8, "bar"))
	assert.False(t, l.ShouldPrintMessage(10, "foo"))
	assert.True(t, l.ShouldPrintMessage(11, "foo"))

	// Verify stale entries are evicted: after t=11, "foo" entry from t=1 is gone.
	// At t=21, "foo" cooldown from t=11 expires; map should only hold active keys.
	assert.True(t, l.ShouldPrintMessage(21, "foo"))
	assert.Equal(t, 1, len(l.nextAllowed)) // only "foo" active; "bar" evicted at t=21
}
