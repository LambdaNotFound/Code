package design

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cell_SetContent(t *testing.T) {
	c := &cell{alias: "A1"}
	c.SetContent("=SUM(B1,C1)")
	assert.Equal(t, "=SUM(B1,C1)", c.content)

	c.SetContent("42")
	assert.Equal(t, "42", c.content)
}
