package oodesign

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Command pattern
 *
 * Encapsulates a request as an object with Execute/Undo methods, so an
 * invoker can run, queue, or reverse arbitrary actions without knowing what
 * any individual command actually does. This is what makes a generic undo
 * stack possible: History only ever calls Command methods, never anything
 * specific to inserting or deleting text.
 */

// Command is what every concrete command implements.
type Command interface {
	Execute()
	Undo()
}

// Document is the receiver — the object commands actually act on.
type Document struct {
	text string
}

// InsertCommand appends text to a Document and removes exactly that text
// again on Undo.
type InsertCommand struct {
	doc  *Document
	text string
}

func NewInsertCommand(doc *Document, text string) *InsertCommand {
	return &InsertCommand{doc: doc, text: text}
}

func (c *InsertCommand) Execute() {
	c.doc.text += c.text
}

func (c *InsertCommand) Undo() {
	c.doc.text = strings.TrimSuffix(c.doc.text, c.text)
}

// DeleteCommand removes n characters from the end of a Document and
// remembers what it removed so Undo can restore it.
type DeleteCommand struct {
	doc     *Document
	n       int
	removed string
}

func NewDeleteCommand(doc *Document, n int) *DeleteCommand {
	return &DeleteCommand{doc: doc, n: n}
}

func (c *DeleteCommand) Execute() {
	cut := c.n
	if cut > len(c.doc.text) {
		cut = len(c.doc.text)
	}
	split := len(c.doc.text) - cut
	c.removed = c.doc.text[split:]
	c.doc.text = c.doc.text[:split]
}

func (c *DeleteCommand) Undo() {
	c.doc.text += c.removed
}

// History is the invoker: it executes commands and can undo the most
// recent one, entirely through the Command interface.
type History struct {
	executed []Command
}

func (h *History) Do(cmd Command) {
	cmd.Execute()
	h.executed = append(h.executed, cmd)
}

func (h *History) Undo() bool {
	if len(h.executed) == 0 {
		return false
	}
	last := h.executed[len(h.executed)-1]
	h.executed = h.executed[:len(h.executed)-1]
	last.Undo()
	return true
}

func Test_InsertCommand_executeAppendsText(t *testing.T) {
	doc := &Document{text: "Hello"}
	cmd := NewInsertCommand(doc, " World")

	cmd.Execute()

	assert.Equal(t, "Hello World", doc.text)
}

func Test_InsertCommand_undoRemovesAppendedText(t *testing.T) {
	doc := &Document{text: "Hello"}
	cmd := NewInsertCommand(doc, " World")
	cmd.Execute()

	cmd.Undo()

	assert.Equal(t, "Hello", doc.text)
}

func Test_DeleteCommand_executeRemovesTrailingCharacters(t *testing.T) {
	doc := &Document{text: "Hello World"}
	cmd := NewDeleteCommand(doc, 6)

	cmd.Execute()

	assert.Equal(t, "Hello", doc.text)
}

func Test_DeleteCommand_undoRestoresRemovedText(t *testing.T) {
	doc := &Document{text: "Hello World"}
	cmd := NewDeleteCommand(doc, 6)
	cmd.Execute()

	cmd.Undo()

	assert.Equal(t, "Hello World", doc.text)
}

func Test_History_undoReversesTheMostRecentCommand(t *testing.T) {
	doc := &Document{}
	history := &History{}

	history.Do(NewInsertCommand(doc, "Hello"))
	history.Do(NewInsertCommand(doc, " World"))
	assert.Equal(t, "Hello World", doc.text)

	history.Undo()

	assert.Equal(t, "Hello", doc.text)
}

func Test_History_undoOnEmptyHistoryReturnsFalse(t *testing.T) {
	history := &History{}

	assert.False(t, history.Undo())
}

func Test_History_invokerNeverNamesAConcreteCommandType(t *testing.T) {
	doc := &Document{text: "Hello World"}
	history := &History{}

	// A mix of command types run through the same Do/Undo calls — History
	// never branches on which kind of Command it's holding.
	history.Do(NewDeleteCommand(doc, 6))   // "Hello World" -> "Hello"
	history.Do(NewInsertCommand(doc, "!")) // "Hello" -> "Hello!"
	assert.Equal(t, "Hello!", doc.text)

	history.Undo() // undo insert
	history.Undo() // undo delete

	assert.Equal(t, "Hello World", doc.text)
	assert.False(t, history.Undo()) // nothing left to undo
}
